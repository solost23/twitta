package routers

import (
	"Twitta/pkg/response"
	"Twitta/pkg/utils"
	"Twitta/services"
	"github.com/gin-gonic/gin"
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
