package routers

import (
	"net/http"

	gcasbin "github.com/maxwellhertz/gin-casbin"

	_ "twitta/docs"
	"twitta/pkg/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRouters(r *gin.Engine) {
	// consul健康检查
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	apiGroup := r.Group("api/twitta")
	{
		// 用户注册
		apiGroup.POST("register", register)
		// 上传头像
		apiGroup.POST("register/avatar", uploadAvatar)
		// 用户登录
		apiGroup.POST("login", login)
		// 用户扫脸登录
		apiGroup.POST("face", face)

		// 展示所有推文
		apiGroup.GET("tweets", tweetList)
		// 用户搜索 - 采用全局搜索
		apiGroup.GET("users/search", userSearch)
		// 推文搜索 - 采用全局搜索
		apiGroup.GET("tweets/search", tweetSearch)

		// swagger
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	apiGroup.Use(
		middlewares.JWTAuth(),
		middlewares.NewCasbinMiddleware().RequiresRoles([]string{"admin", "user"}, gcasbin.WithLogic(gcasbin.OR)),
	)
	{
		// 用户相关
		initAuthUserRouter(apiGroup)
		// 推文相关
		initAuthTweetRouter(apiGroup)
		// 交友相关
		initAuthFriendRouter(apiGroup)
		// 聊天相关
		initAuthChatRouter(apiGroup)
		// 关注-粉丝相关
		initFansRouter(apiGroup)
		// 点赞-评论相关
		initCommentRouter(apiGroup)
	}
}

func initAuthUserRouter(group *gin.RouterGroup) {
	user := group.Group("users")
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
	tweet := group.Group("tweets")
	{
		// 发送推文
		tweet.POST("", tweetSend)
		// 推文文件上传
		tweet.POST("static", staticUpload)
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
	friend := group.Group("friends")
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
	chat := group.Group("chats")
	{
		// 聊天信息列表
		chat.GET(":id", chatList)
	}
}

func initFansRouter(group *gin.RouterGroup) {
	fan := group.Group("fans")
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
	comment := group.Group("comments")
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
