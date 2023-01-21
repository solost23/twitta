package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/solost23/go_interface/gen_go/common"
	"github.com/solost23/go_interface/gen_go/push"
	"mime/multipart"
	"time"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/cache"
	"twitta/pkg/constants"
	"twitta/pkg/middlewares"
	"twitta/pkg/models"
	"twitta/pkg/utils"

	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func (s *Service) Register(c *gin.Context, params *forms.RegisterForm) error {
	db := global.DB

	user := models.NewUser()
	err := user.FindOne(c, db, constants.Mongo, bson.M{"username": params.Username}, user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("用户已存在，不允许重复创建")
	}
	data := &models.User{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:        utils.UUID(),
		Username:  params.Username,
		Password:  utils.NewMd5(params.Password, constants.Secret),
		Nickname:  params.Nickname,
		Mobile:    params.Mobile,
		Role:      "user",
		Avatar:    utils.TrimDomainPrefix(params.Avatar),
		Introduce: params.Introduce,
		Email:     params.Email,
		FansCount: 0,
		Disabled:  0,
	}
	_, err = user.InsertOne(c, db, constants.Mongo, data)
	if err != nil {
		return err
	}

	// 将用户数据存入zinc
	err = NewZinc().InsertDocument(c, constants.ZINCINDEXUSER, data.ID, map[string]interface{}{
		"basemodel": map[string]interface{}{
			"created-at": data.BaseModel.CreatedAt,
			"updated-at": data.BaseModel.UpdatedAt,
			"deleted-at": data.BaseModel.DeletedAt,
		},
		"username":        data.Username,
		"password":        data.Password,
		"nickname":        data.Nickname,
		"mobile":          data.Mobile,
		"role":            data.Role,
		"avatar":          utils.TrimDomainPrefix(data.Avatar),
		"introduce":       data.Introduce,
		"email":           data.Email,
		"fans_count":      data.FansCount,
		"wechat_count":    data.WechatCount,
		"last_login_time": data.LastLoginTime,
		"disabled":        data.Disabled,
	})
	if err != nil {
		return err
	}

	// 调用邮件发送服务发送邮件
	if len(params.Email) >= 0 {
		reply, err := global.PushSrvClient.SendEmail(c, &push.SendEmailRequest{
			Header: &common.RequestHeader{
				TraceId:     6678677,
				OperatorUid: 55,
			},
			Email: &push.Email{
				Topic:       "register",
				Name:        params.Username,
				Addr:        params.Email,
				ContentType: "text/plain",
				Content:     fmt.Sprintf("恭喜%s注册Twitta成功", params.Username),
			},
		})
		if err != nil {
			return err
		}
		if reply.ErrorInfo.GetCode() != 0 {
			return errors.New(reply.ErrorInfo.GetMsg())
		}
	}
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
	db := global.DB

	user := &models.User{}
	err := user.FindOne(c, db, constants.Mongo, bson.M{"username": params.Username}, user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if params.Username != user.Username || utils.NewMd5(params.Password, constants.Secret) != user.Password {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Disabled == 1 {
		return nil, errors.New("您的账户已被禁用，请联系管理员")
	}
	// 区分两种设备 分别是web和mobile
	var redisPrefix string
	if params.Device == "web" {
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
		return nil, err
	}
	userJson, _ := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.RedisConnFactory(10)
	if err != nil {
		return nil, err
	}

	key := redisPrefix + user.ID
	oldToken, err := rdb.Get(c, key).Result()

	rdb.Del(c, constants.RedisPrefix+oldToken)
	rdb.Set(c, key, token, time.Duration(global.ServerConfig.JWTConfig.Duration)*time.Second)
	rdb.Set(c, constants.RedisPrefix+token, userJson, time.Duration(global.ServerConfig.JWTConfig.Duration)*time.Second)

	_, err = user.Update(c, db, constants.Mongo, bson.M{"_id": user.ID}, bson.D{{"$set", bson.D{{"last_login_time", time.Now()}}}})
	if err != nil {
		return nil, err
	}

	response := &forms.LoginResponse{
		IsFirstLogin: 2,
		User:         *user,
		Token:        token,
	}

	// 调用邮件发送服务发送邮件
	if len(user.Email) >= 0 {
		reply, err := global.PushSrvClient.SendEmail(c, &push.SendEmailRequest{
			Header: &common.RequestHeader{
				TraceId:     6678678,
				OperatorUid: 56,
			},
			Email: &push.Email{
				Topic:       "login",
				Name:        user.Username,
				Addr:        user.Email,
				ContentType: "text/plain",
				Content:     fmt.Sprintf("恭喜%s登陆Twitta成功", user.Username),
			},
		})
		if err != nil {
			return nil, err
		}
		if reply.ErrorInfo.GetCode() != 0 {
			return nil, errors.New(reply.ErrorInfo.GetMsg())
		}
	}

	return response, err
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
	db := global.DB
	user := utils.GetUser(c)

	update := bson.M{
		"$set": bson.M{
			"username":  params.Username,
			"nickname":  params.Nickname,
			"avatar":    utils.TrimDomainPrefix(params.Avatar),
			"introduce": params.Introduce,
		},
	}
	_, err := models.NewUser().Update(c, db, constants.Mongo, bson.M{"_id": user.ID}, update)
	if err != nil {
		return err
	}

	data := &models.User{}
	err = models.NewUser().FindOne(c, db, constants.Mongo, bson.M{"_id": user.ID}, &data)
	if err != nil {
		return err
	}
	// 拿到id 更新zinc数据
	// 删除 + 插入 = 更新
	err = NewZinc().DeleteDocument(c, constants.ZINCINDEXUSER, user.ID)
	if err != nil {
		return err
	}
	err = NewZinc().InsertDocument(c, constants.ZINCINDEXUSER, user.ID, map[string]interface{}{
		"basemodel": map[string]interface{}{
			"created-at": data.BaseModel.CreatedAt,
			"updated-at": data.BaseModel.UpdatedAt,
			"deleted-at": data.BaseModel.DeletedAt,
		},
		"username":        data.Username,
		"password":        data.Password,
		"nickname":        data.Nickname,
		"mobile":          data.Mobile,
		"role":            data.Role,
		"avatar":          utils.TrimDomainPrefix(data.Avatar),
		"introduce":       data.Introduce,
		"email":           data.Email,
		"fans_count":      data.FansCount,
		"wechat_count":    data.WechatCount,
		"last_login_time": data.LastLoginTime,
		"disabled":        data.Disabled,
	})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) UserDetail(c *gin.Context, id string) (*forms.UserDetailResponse, error) {
	db := global.DB

	user := &models.User{}
	err := models.NewUser().FindOne(c, db, constants.Mongo, bson.M{"_id": id}, &user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New(fmt.Sprintf("不存在此用户"))
	}
	userDetailResponse := &forms.UserDetailResponse{
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

func (*Service) UserSearch(c *gin.Context, params *forms.SearchForm) ([]*forms.UserDetailResponse, error) {
	db := global.DB

	// 直接从zinc中搜索数据，然后返回搜索到的数据
	from := int32((params.Page - 1) * params.Size)
	size := from + int32(params.Size) - 1
	searchResults, _, err := NewZinc().SearchDocument(c, constants.ZINCINDEXUSER, params.Keyword, from, size)
	if err != nil {
		return nil, err
	}

	// 求出粉丝数和关注数的映射关系
	userIds := make([]string, 0, len(searchResults))
	for _, searchResult := range searchResults {
		userIds = append(userIds, *searchResult.Id)
	}
	users := make([]*models.User, 0, len(searchResults))
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	if err != nil {
		return nil, err
	}
	userIdToFansWechatNumMaps := make(map[string]struct {
		FansCount   int64
		WechatCount int64
	}, len(users))
	for _, user := range users {
		userIdToFansWechatNumMaps[user.ID] = struct {
			FansCount   int64
			WechatCount int64
		}{FansCount: user.FansCount, WechatCount: user.WechatCount}
	}
	userSearchResponse := make([]*forms.UserDetailResponse, 0, len(searchResults))
	for _, searchResult := range searchResults {
		userSearchResponse = append(userSearchResponse, &forms.UserDetailResponse{
			UserId:      *searchResult.Id,
			Username:    searchResult.Source["username"].(string),
			Nickname:    searchResult.Source["nickname"].(string),
			Avatar:      searchResult.Source["avatar"].(string),
			Introduce:   searchResult.Source["introduce"].(string),
			WechatCount: userIdToFansWechatNumMaps[*searchResult.Id].WechatCount,
			FansCount:   userIdToFansWechatNumMaps[*searchResult.Id].FansCount,
			CreatedAt:   searchResult.Source["basemodel"].(map[string]interface{})["created-at"].(string),
		})
	}
	return userSearchResponse, nil
}
