package routers

import (
	"github.com/gin-gonic/gin"
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"
)

func chatList(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	result, err := services.NewService().ChatList(c, UIdForm.Id)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}
