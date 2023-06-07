package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	servantEs "twitta/services/servants/es"
	servantPush "twitta/services/servants/push"
	servantUser "twitta/services/servants/user"

	"github.com/solost23/protopb/gen/go/protos/common"
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"go.uber.org/zap"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/cache"
	"twitta/pkg/constants"
	"twitta/pkg/domain"
	"twitta/pkg/middlewares"
	"twitta/pkg/models"
	"twitta/pkg/utils"

	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func (s *Service) Register(c *gin.Context, params *forms.RegisterForm) error {
	db := models.NewDB()

	_, err := models.GWhereFirst[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"username": params.Username})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("用户已存在，不允许重复创建")
	}

	// 获取用户头像编码
	faceImg := params.FaceImg
	faceEncoding := ""
	if faceImg != "" {
		resp, err := http.Get(faceImg)
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		faceEncodings, err := servantUser.GenerateFaceEncodings(c, [][]byte{bytes})
		if err != nil {
			return err
		}
		if len(faceEncodings) > 0 {
			faceEncoding = faceEncodings[0]
		}
	}

	data := &models.User{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:           utils.UUID(),
		Username:     params.Username,
		Password:     utils.NewMd5(params.Password, constants.Secret),
		Nickname:     params.Nickname,
		Mobile:       params.Mobile,
		Role:         "user",
		Avatar:       utils.TrimDomainPrefix(params.Avatar),
		Introduce:    params.Introduce,
		Email:        params.Email,
		FansCount:    0,
		Disabled:     0,
		FaceImg:      utils.TrimDomainPrefix(faceImg),
		FaceEncoding: faceEncoding,
	}
	_, err = models.GInsertOne[models.User](c, db.GetCollection(models.NewUser().TableName()), data)
	if err != nil {
		return err
	}

	// 用户数据存入es
	go func() {
		if err = servantEs.Save(c, constants.ESCINDEXUSER, data.ID, data); err != nil {
			zap.S().Error(err)
		}
	}()

	// 发送邮件
	go func() {
		if err = servantPush.SendEmail(c, "注册", params.Username, params.Email, fmt.Sprintf("恭喜%s注册Twitta成功", params.Username)); err != nil {
			zap.S().Error(err)
		}
	}()

	return nil
}

func (s *Service) UploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {
	folder := "twitta.users.avatar"

	url, err := UploadImg(0, folder, file.Filename, file)
	if err != nil {
		return "", err
	}
	// 对链接做处理
	// eg:http://minio:9000/avatar/5ac8dd9f599264da59532bc31593b7b7.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20221125%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20221125T033740Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=514d733ebc7a4e1b58d54b810a191dd9994d8745231045c7a0908b895c0c14db
	// 应处理成: http://localhost:9000/avatar/5ac8dd9f599264da59532bc31593b7b7.jpeg，返回
	return utils.FulfillImageOSSPrefix(utils.TrimDomainPrefix(url)), nil
}

func (s *Service) Login(c *gin.Context, params *forms.LoginForm) (*forms.LoginResponse, error) {
	db := models.NewDB()

	user, err := models.GWhereFirst[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"username": params.Username})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if params.Username != user.Username || utils.NewMd5(params.Password, constants.Secret) != user.Password {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := loginAndGetToken(c, params.Device, user)
	if err != nil {
		return nil, err
	}

	response := &forms.LoginResponse{
		IsFirstLogin: 2,
		User:         *user,
		Token:        token,
	}

	// 发送邮件
	go func() {
		if err = servantPush.SendEmail(c, "登录", user.Username, user.Email, fmt.Sprintf("恭喜%s登录Twitta成功", user.Username)); err != nil {
			zap.S().Error(err)
		}
	}()

	return response, err
}

func (s *Service) Face(c *gin.Context, file *multipart.FileHeader) (*forms.Face, error) {
	fp, err := file.Open()
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	userId, isFound, err := servantUser.CompareFace(c, b)
	if err != nil {
		return nil, err
	}
	result := &forms.Face{
		IsFound: isFound,
	}
	if !isFound {
		return result, nil
	}

	db := models.NewDB()

	user, err := models.GWhereFirst[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"_id": userId})
	if err != nil {
		return nil, err
	}
	token, err := loginAndGetToken(c, "web", user)
	if err != nil {
		return nil, err
	}
	result.User = *user
	result.Token = token

	// 发送邮件
	go func() {
		if err = servantPush.SendEmail(c, "登录", user.Username, user.Email, fmt.Sprintf("恭喜%s登录Twitta成功", user.Username)); err != nil {
			zap.S().Error(err)
		}
	}()

	return result, nil
}

func loginAndGetToken(ctx context.Context, device string, user *models.User) (string, error) {
	if user.Disabled == 1 {
		return "", errors.New("您的账户已被禁用，请联系管理员")
	}

	// 区分两种设备 分别是web和mobile
	var redisPrefix string
	if device == "web" {
		redisPrefix = constants.WebRedisPrefix
	} else {
		redisPrefix = constants.MobileRedisPrefix
	}

	j := middlewares.NewJWT()
	claims := middlewares.CustomClaims{
		UserId: user.ID,
		Device: redisPrefix,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + global.ServerConfig.JWTConfig.Duration,
			Issuer:    "twitta",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "nil", err
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return "nil", err
	}
	rdb, err := cache.RedisConnFactory(10)
	if err != nil {
		return "", err
	}

	key := redisPrefix + user.ID
	oldToken, err := rdb.Get(ctx, key).Result()

	rdb.Del(ctx, constants.RedisPrefix+oldToken)
	rdb.Set(ctx, key, token, time.Duration(global.ServerConfig.JWTConfig.Duration)*time.Second)
	rdb.Set(ctx, constants.RedisPrefix+token, userJson, time.Duration(global.ServerConfig.JWTConfig.Duration)*time.Second)

	_, err = models.GWhereUpdate[models.User](ctx, models.NewDB().GetCollection(models.NewUser().TableName()), bson.M{"_id": user.ID}, bson.D{{"$set", bson.D{{"last_login_time", time.Now()}}}})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Logout(c *gin.Context, params *forms.LogoutForm) error {
	user := utils.GetUser(c)

	rdb, err := cache.RedisConnFactory(10)
	if err != nil {
		return err
	}
	var redisPrefix string
	if params.Device == "web" {
		redisPrefix = constants.WebRedisPrefix
	} else {
		redisPrefix = constants.MobileRedisPrefix
	}
	key := redisPrefix + user.ID
	token, err := rdb.Get(c, key).Result()
	if err != nil {
		return err
	}
	rdb.Del(c, constants.RedisPrefix+token)
	return nil
}

func (*Service) UserUpdate(c *gin.Context, params *forms.UserUpdateForm) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	data := bson.M{}
	if params.Username != "" {
		data["username"] = params.Username
	}
	if params.Nickname != "" {
		data["nickname"] = bson.M{}
	}
	if params.Avatar != "" {
		data["avatar"] = utils.TrimDomainPrefix(params.Avatar)
	}
	if params.Introduce != "" {
		data["introduce"] = params.Introduce
	}
	// 更新用户头像编码
	if params.FaceImg != "" {
		resp, err := http.Get(params.FaceImg)
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		faceEncodings, err := servantUser.GenerateFaceEncodings(c, [][]byte{bytes})
		if err != nil {
			return err
		}
		if len(faceEncodings) > 0 {
			data["face_img"] = utils.TrimDomainPrefix(params.FaceImg)
			data["face_encoding"] = faceEncodings[0]
		}
	}
	update := bson.M{"$set": data}
	_, err := models.GWhereUpdate[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"_id": user.ID}, update)
	if err != nil {
		return err
	}
	user, err = models.GWhereFirst[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"_id": user.ID})
	if err != nil {
		return err
	}

	// 拿到id 更新es数据
	// 删除 + 插入 = 更新
	go func() {
		if err = servantEs.Delete(c, constants.ESCINDEXUSER, user.ID); err != nil {
			zap.S().Error(err)
		}
		if err = servantEs.Save(c, constants.ESCINDEXUSER, user.ID, user); err != nil {
			zap.S().Error(err)
		}
	}()

	return nil
}

func (*Service) UserDetail(c *gin.Context, id string) (*forms.UserDetail, error) {
	db := models.NewDB()

	user, err := models.GWhereFirst[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New(fmt.Sprintf("不存在此用户"))
	}
	userDetailResponse := &forms.UserDetail{
		UserId:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      utils.FulfillImageOSSPrefix(user.Avatar),
		Introduce:   user.Introduce,
		WechatCount: user.WechatCount,
		FansCount:   user.FansCount,
		CreatedAt:   user.CreatedAt.Format(constants.TimeFormat),
	}
	return userDetailResponse, nil
}

func (*Service) UserSearch(c *gin.Context, params *forms.SearchForm) (*forms.UserSearch, error) {

	searchResult, err := domain.NewESClient().Search(c, &es_service.SearchRequest{
		Header: &common.RequestHeader{
			Requester:   "search_user",
			OperatorUid: -1,
		},
		ShouldQuery: &es_service.Query{
			MultiMatchQueries: []*es_service.MultiMatchQuery{
				{Field: []string{"username", "nickname", "mobile", "role", "introduce", "email"}, Value: params.Keyword},
			},
		},
		Indices: []string{constants.ESCINDEXUSER},
		Page:    int32(params.Page),
		Size:    int32(params.Size),
		Pretty:  true,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*forms.UserDetail, 0, len(searchResult.Records))
	for _, search := range searchResult.Records {
		record := &forms.UserDetail{}
		_ = json.Unmarshal([]byte(search), record)
		records = append(records, record)
	}

	result := &forms.UserSearch{
		Records: records,
		PageList: utils.PageList{
			Size:    int(searchResult.PageList.GetSize()),
			Pages:   searchResult.PageList.GetPages(),
			Total:   searchResult.PageList.GetTotal(),
			Current: int(searchResult.PageList.GetCurrent()),
		},
	}
	return result, nil
}
