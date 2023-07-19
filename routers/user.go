package routers

import (
	"twitta/forms"
	"twitta/pkg/utils"
	"twitta/services"

	"twitta/pkg/response"

	"github.com/gin-gonic/gin"
)

//@Summary register
//@Tags user
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param register body forms.RegisterForm true "registerForm"
//@Router /register [post]
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

//@Summary upload avatar
//@Tags user
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param file formData file true "userAvatar"
//@Router /register/avatar [post]
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

//	@Summary	login
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		login	body		forms.LoginForm	true	"Login credentials"
//	@Success	200		{object}	forms.LoginResponse
//	@Failure	400		{object}	response.Response
//	@Router		/login [post]
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

func face(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	result, err := services.NewService().Face(c, file)
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

//@Summary user info update
//@Tags user
//@Produce json
//@Success 200 {object} response.Response
//@Failure 400 {object} response.Response
//@Param userUpdateForm body forms.UserUpdateForm true "userUpdateForm"
//	@Param	token	header	string	true	"token"
//@Router /users [put]
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

//@Summary user detail
//@Tags user
//@Produce json
//@Success 200 {object} forms.UserDetail
//@Failure 400 {object} response.Response
//@Param id path string true "userId"
//	@Param	token	header	string	true	"token"
//@Router /users/{id} [get]
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

//@Summary search user
//@Tags user
//@Produce json
//@Success 200 {object} forms.UserSearch
//@Failure 400 {object} response.Response
//@Param searchForm body forms.SearchForm true "searchForm"
//@Router /users/search [get]
func userSearch(c *gin.Context) {
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
	result, err := services.NewService().UserSearch(c, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}

	response.Success(c, result)
}
