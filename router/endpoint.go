package router

import (
	"Twitta/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	group := router.Group("api")
	initNoAuthRouter(group)
	group.Use(middlewares.JWTAuth())
	initAuthRouter(group)
	return router
}

func initNoAuthRouter(group *gin.RouterGroup) {
	group.POST("register", register)
	group.POST("register/avatar", uploadAvatar)
	group.POST("login", login)
}

func initAuthRouter(group *gin.RouterGroup) {
	initAuthUserRouter(group)
}

func initAuthUserRouter(group *gin.RouterGroup) {
	user := group.Group("user")
	{
		// 注销用户
		user.POST("logout", logout)
	}
}
