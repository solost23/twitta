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
	params := &utils.PageForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}

	result, err := services.NewService().ChatList(c, UIdForm.Id, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}
