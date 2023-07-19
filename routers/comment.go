package routers

import (
	"twitta/forms"
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"

	"github.com/gin-gonic/gin"
)

//	@Summary	commentList
//	@Tags		thumb
//	@Produce	json
//	@Success	200				{object}	forms.CommentList
//	@Failure	400				{object}	response.Response
//	@Param		id				path		string	true	"commentId"
//	@Router		/comments/{id}	[get]
func commentList(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	params := &forms.CommentInsertForm{}
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

	result, err := services.NewService().CommentList(c, UIdForm.Id, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//	@Summary	create thumb
//	@Tags		thumb
//	@Produce	json
//	@Success	200	{object}	response.Response
//	@Failure	400	{object}	response.Response
//	@Param		id	path		string	true	"commentId"
//	@Router		/comments/{id}/thumb [post]
func commentThumb(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().CommentThumb(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//	@Summary	delete thumb
//	@Tags		thumb
//	@Produce	json
//	@Success	200	{object}	response.Response
//	@Failure	400	{object}	response.Response
//	@Param		id	path		string	true	"commentId"
//	@Router		/comments/{id}/thumb [delete]
func commentThumbDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().CommentThumbDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//	@Summary	create comment
//	@Tags		comment
//	@Produce	json
//	@Success	200				{object}	response.Response
//	@Failure	400				{object}	response.Response
//	@Param		id				path		string					true	"tweetId"
//	@Param		commentInsert	body		forms.CommentInsertForm	true	"commentInsertForm"
//	@Router		/comments/{id} [post]
func commentInsert(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	params := &forms.CommentInsertForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().CommentInsert(c, UIdForm.Id, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//	@Summary	delete comment
//	@Tags		comment
//	@Produce	json
//	@Success	200	{object}	response.Response
//	@Failure	400	{object}	response.Response
//	@Param		id	path		string	true	"commentId"
//	@Router		/comments/{id} [delete]
func commentDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().CommentDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}
