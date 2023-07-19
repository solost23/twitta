package routers

import (
	"twitta/forms"
	"twitta/pkg/constants"
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/services"

	"github.com/gin-gonic/gin"
)

//@Summary send tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param tweetCreateForm body forms.TweetCreateForm true "tweetCreateForm"
//@Router /tweets [post]
func tweetSend(c *gin.Context) {
	params := &forms.TweetCreateForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetSend(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//@Summary tweet file upload
//@Tags tweet
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param file formData file true "file"
//@Router /tweets/static [post]
func staticUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, constants.BadRequestCode, err)
		return
	}
	result, err := services.NewService().StaticUpload(c, file)
	if err != nil {
		response.Error(c, constants.InternalServerErrorCode, err)
		return
	}

	response.Success(c, result)
}

//@Summary delete tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "tweetId"
//@Router /tweets/{id} [delete]
func tweetDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//@Summary tweet list
//@Tags tweet
//@Produce json
//@Success 200 {object} forms.TweetList
//@Failure 400 {object} response.Response
//@Param pageForm body utils.PageForm true "pageForm"
//@Router /tweets [get]
func tweetList(c *gin.Context) {
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
	result, err := services.NewService().TweetList(c, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary favorite tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param tweetFavoriteForm body forms.TweetFavoriteForm true "tweetFavoriteForm"
//@Router /tweets/favorite [post]
func tweetFavorite(c *gin.Context) {
	params := &forms.TweetFavoriteForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetFavorite(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//@Summary cancel favorite tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param id path string true "tweetId"
//@Router /tweets/favorite/{id} [delete]
func tweetFavoriteDelete(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	if err := services.NewService().TweetFavoriteDelete(c, UIdForm.Id); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.MessageSuccess(c, "成功", nil)
}

//@Summary favorite tweet list
//@Tags tweet
//@Produce json
//@Success 200 {object} forms.TweetList
//@Failure 400 {object} response.Response
//@Router /tweets/favorite [get]
func tweetFavoriteList(c *gin.Context) {
	result, err := services.NewService().TweetFavoriteList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary user tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} forms.TweetList
//@Failure 400 {object} response.Response
//@Router /tweets/own [get]
func tweetOwnList(c *gin.Context) {
	result, err := services.NewService().TweetOwnList(c)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//@Summary search tweet
//@Tags tweet
//@Produce json
//@Success 200 {object} forms.TweetList
//@Failure 400 {object} response.Response
//@Param searchForm body forms.SearchForm true "searchForm"
//@Router /tweets/search [get]
func tweetSearch(c *gin.Context) {
	params := &forms.SearchForm{}
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

	result, err := services.NewService().TweetSearch(c, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}

	response.Success(c, result)
}
