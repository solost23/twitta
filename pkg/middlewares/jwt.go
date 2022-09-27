package middlewares

import (
	"Twitta/global"
	"Twitta/pkg/cache"
	"Twitta/pkg/constants"
	"Twitta/pkg/models"
	"Twitta/pkg/response"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.Error(c, 1999, errors.New("请登录"))
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				response.Error(c, 1998, errors.New("授权已过期"))
				return
			}
			response.Error(c, 1999, errors.New("无效 token"))
			return
		}
		rdb, err := cache.RedisConnFactory(10)
		if err != nil {
			response.Error(c, 1999, err)
			return
		}
		jwtToken, err := rdb.Get(c, claims.Device+claims.UserId).Result()
		if err != nil {
			response.Error(c, 1999, errors.New("无效 token"))
			return
		}
		if jwtToken != token {
			response.Error(c, 1999, errors.New("无效 token"))
			return
		}
		jsonUser, err := rdb.Get(c, constants.RedisPrefix+token).Result()
		if err != nil {
			response.Error(c, 1999, errors.New("无效 token"))
			return
		}

		redisUser := &models.User{}
		err = json.Unmarshal([]byte(jsonUser), redisUser)
		if err != nil {
			response.Error(c, 1999, err)
			return
		}
		c.Set("user", redisUser)
		c.Next()
		return
	}
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	UserId string
	Device string
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{[]byte(global.ServerConfig.JWTConfig.Key)}
}

// CreateToken 创建 Token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Toke is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// RefreshToken 更新 Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(global.ServerConfig.JWTConfig.Duration) * time.Second).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
