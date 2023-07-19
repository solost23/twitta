package routers

import (
	"twitta/forms"
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"

	"github.com/gin-gonic/gin"
)

//@Summary send friend application
//@Tags friend
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Router /friends [delete]
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

//@Summary accept friend application
//@Tags friend
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "userId"
//@Router /friends/{id}/accept [put]
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

//@Summary reject friend application
//@Tags friend
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "userId"
//@Router /friends/{id}/reject [put]
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

//@Summary application list
//@Tags friend
//@Produce json
//@Success 200 {object} forms.FriendApplicationListResponse
//@Failure 400 {object} response.Response
//@Router /friends [get]
func friendApplicationList(c *gin.Context) {
	result, err := services.NewService().FriendApplicationList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary delete friend
//@Tags friend
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "userId"
//@Router /friends/{id} [delete]
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
