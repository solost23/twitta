package routers

import (
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"

	"github.com/gin-gonic/gin"
)

//	@Summary	fan list
//	@Tags		fan
//	@Produce	json
//	@Success	200		{object}	forms.FansAndWhatResponse
//	@Failure	400		{object}	response.Response
//	@Param		token	header		string	true	"token"
//	@Router		/fans [get]
func fanList(c *gin.Context) {
	result, err := services.NewService().FanList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary what list
//@Tags what
//@Produce json
//@Success 200 {object} forms.FansAndWhatResponse
//@Failure 400 {object} response.Response
//	@Param	token	header	string	true	"token"
//@Router /fans/what [get]
func whatList(c *gin.Context) {
	result, err := services.NewService().WhatList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary what a user
//@Tags what
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "whatUserId"
//	@Param	token	header	string	true	"token"
//@Router /fans/{id} [post]
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

//@Summary cancel what
//@Tags what
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "whatUserId"
//	@Param	token	header	string	true	"token"
//@Router /fans/{id} [delete]
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
