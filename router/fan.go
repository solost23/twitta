package router

import (
	"Twitta/pkg/response"
	"Twitta/pkg/utils"
	"Twitta/services"
	"github.com/gin-gonic/gin"
)

func fanList(c *gin.Context) {
	result, err := services.NewService().FanList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

func whatList(c *gin.Context) {
	result, err := services.NewService().WhatList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

func whatUser(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().WhatUser(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

func whatUserDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().WhatUserDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}
