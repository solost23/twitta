package router

import (
	"Twitta/forms"
	"Twitta/pkg/utils"
	"Twitta/services"

	"Twitta/pkg/response"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	params := &forms.RegisterForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().Register(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func uploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	result, err := services.NewService().UploadAvatar(c, file)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

func login(c *gin.Context) {
	params := &forms.LoginForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	result, err := services.NewService().Login(c, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

func logout(c *gin.Context) {
	params := &forms.LogoutForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().Logout(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func userUpdate(c *gin.Context) {
	params := &forms.UserUpdateForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().UserUpdate(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func userDetail(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	result, err := services.NewService().UserDetail(c, UIdForm.Id)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}

	response.Success(c, result)
}
