package services

import (
	"Twitta/forms"
	"Twitta/global"
	"Twitta/pkg/cache"
	"Twitta/pkg/constants"
	"Twitta/pkg/middlewares"
	"Twitta/pkg/models"
	"Twitta/pkg/utils"
	"encoding/json"
	"errors"
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func (s *Service) Register(c *gin.Context, params *forms.RegisterForm) error {
	db := global.DB

	user := models.NewUser()
	err := user.FindOne(c, db, constants.Mongo, user.TableName(), bson.M{"username": params.Username}, user)
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
		Avatar:    params.Avatar,
		Introduce: params.Introduce,
		FansCount: 0,
		Disabled:  0,
	}
	_, err = user.InsertOne(c, db, constants.Mongo, user.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {
	user := &models.User{}
	result, err := UploadImg(user, "avatar", file)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *Service) Login(c *gin.Context, params *forms.LoginForm) (*forms.LoginResponse, error) {
	db := global.DB

	user := &models.User{}
	err := user.FindOne(c, db, constants.Mongo, user.TableName(), bson.M{"username": params.Username}, user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if params.Username != user.Username || utils.NewMd5(params.Password, constants.Secret) != user.Password {
		return nil, errors.New("用户名或密码错误")
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
			Issuer:    "Twitta",
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

	_, err = user.Update(c, db, constants.Mongo, user.TableName(), bson.M{"_id": user.ID}, bson.D{{"$set", bson.D{{"last-login-time", time.Now()}}}})
	if err != nil {
		return nil, err
	}

	response := &forms.LoginResponse{
		IsFirstLogin: 2,
		User:         *user,
		Token:        token,
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
