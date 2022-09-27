package router

import (
	"Twitta/forms"
	"Twitta/pkg/response"
	"Twitta/pkg/utils"
	"Twitta/services"

	"github.com/gin-gonic/gin"
)

func tweetSend(c *gin.Context) {
	params := &forms.TweetCreateForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetSend(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func tweetDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func tweetList(c *gin.Context) {
	result, err := services.NewService().TweetList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}
