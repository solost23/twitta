package routers

import (
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"

	"github.com/gin-gonic/gin"
)

//	@Summary	chatList
//	@Tags		chat
//	@Produce	json
//	@Success	200		{object}	forms.ChatList
//	@Failure	400		{object}	response.Response
//	@Param		id		path		string	true	"chatId"
//	@Param		token	header		string	true	"token"
//	@Router		/chats/{id} [get]
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
