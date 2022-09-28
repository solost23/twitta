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

	// 展示所有推文
	group.GET("tweet", tweetList)
}

func initAuthRouter(group *gin.RouterGroup) {
	// 用户相关
	initAuthUserRouter(group)
	// 推文相关
	initAuthTweetRouter(group)
	// 交友相关
	initAuthFriendRouter(group)
}

func initAuthUserRouter(group *gin.RouterGroup) {
	user := group.Group("user")
	{
		// 注销用户
		user.POST("logout", logout)
	}
}

func initAuthTweetRouter(group *gin.RouterGroup) {
	tweet := group.Group("tweet")
	{
		// 发送推文
		tweet.POST("", tweetSend)
		// 删除推文
		tweet.DELETE(":id", tweetDelete)
		// 收藏推文
		tweet.POST("favorite", tweetFavorite)
		// 取消收藏推文
		tweet.DELETE("favorite/:id", tweetFavoriteDelete)
		// 展示当前用户收藏推文
		tweet.GET("favorite", tweetFavoriteList)
		// 展示当前用户的推文
		tweet.GET("own", tweetOwnList)
	}
}

func initAuthFriendRouter(group *gin.RouterGroup) {
	friend := group.Group("friend")
	{
		// 发送好友申请
		friend.POST("", friendApplicationSend)
		// 通过好友申请
		friend.PUT(":id/accept", friendApplicationAccept)
		// 拒绝好友申请
		friend.PUT(":id/reject", friendApplicationReject)
		// 好友申请列表
		friend.GET("", friendApplicationList)
		// 删除好友
		friend.DELETE(":id", friendDelete)
	}
}
