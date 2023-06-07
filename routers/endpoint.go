package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "twitta/docs"
	"twitta/pkg/middlewares"
)

func SetRouters(r *gin.Engine) {
	group := r.Group("api/twitta")
	initNoAuthRouter(group)
	group.Use(
		middlewares.JWTAuth(),
		// middlewares.AuthCheckRole(),
	)

	initAuthRouter(group)
}

func initNoAuthRouter(group *gin.RouterGroup) {
	group.POST("register", register)
	group.POST("register/avatar", uploadAvatar)
	group.POST("login", login)
	// 验证用户脸部
	group.POST("face", face)

	// 展示所有推文
	group.GET("tweet", tweetList)
	// 用户搜索 - 采用全局搜索
	group.GET("user/search", userSearch)
	// 推文搜索 - 采用全局搜索
	group.GET("tweet/search", tweetSearch)

	// swagger
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initAuthRouter(group *gin.RouterGroup) {
	// 用户相关
	initAuthUserRouter(group)
	// 推文相关
	initAuthTweetRouter(group)
	// 交友相关
	initAuthFriendRouter(group)
	// 聊天相关
	initAuthChatRouter(group)
	// 关注-粉丝相关
	initFansRouter(group)
	// 点赞-评论相关
	initCommentRouter(group)
}

func initAuthUserRouter(group *gin.RouterGroup) {
	user := group.Group("user")
	{
		// 注销用户
		user.POST("logout", logout)
		// 编辑个人资料(自己用)
		user.PUT("", userUpdate)
		// 展示用户资料（自己和他人都可以用）
		user.GET(":id", userDetail)
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
		// 发送好友申请-朋友私信发送
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

func initAuthChatRouter(group *gin.RouterGroup) {
	chat := group.Group("chat")
	{
		// 聊天信息列表
		chat.GET(":id", chatList)
	}
}

func initFansRouter(group *gin.RouterGroup) {
	fan := group.Group("fan")
	{
		// 个人粉丝列表
		fan.GET("", fanList)
		// 个人关注列表
		fan.GET("what", whatList)
		// 关注某人
		fan.POST(":id", whatUser)
		// 取消关注
		fan.DELETE(":id", whatUserDelete)
	}
}

func initCommentRouter(group *gin.RouterGroup) {
	comment := group.Group("comment")
	{
		// 推文评论列表
		comment.GET(":id", commentList)
		// 用户点赞推文
		comment.POST(":id/thumb", commentThumb)
		// 用户取消点赞推文
		comment.DELETE(":id/thumb", commentThumbDelete)
		// 用户评论推文
		comment.POST(":id", commentInsert)
		// 用户删除评论
		comment.DELETE(":id", commentDelete)
	}
}
