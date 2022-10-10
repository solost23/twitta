package routers

import (
	"Twitta/forms"
	"Twitta/pkg/response"
	"Twitta/pkg/utils"
	"Twitta/services"
	"github.com/gin-gonic/gin"
)

func friendApplicationSend(c *gin.Context) {
	params := &forms.FriendApplicationSendForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	err := services.NewService().FriendApplicationSend(c, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func friendApplicationAccept(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	err := services.NewService().FriendApplicationAccept(c, UIdForm.Id)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func friendApplicationReject(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	err := services.NewService().FriendApplicationReject(c, UIdForm.Id)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func friendApplicationList(c *gin.Context) {
	result, err := services.NewService().FriendApplicationList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

func friendDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	err := services.NewService().FriendDelete(c, UIdForm.Id)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}
