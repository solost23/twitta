package forms

import (
	"twitta/pkg/models"
	"twitta/pkg/utils"
)

type RegisterForm struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Nickname  string `json:"nickname"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email" binding:"required"`
	Avatar    string `json:"avatar"`
	Introduce string `json:"introduce"`
	FaceImg   string `json:"faceImg"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Platform string `json:"platform" comment:"平台" binding:"required,oneof=twitta video_server"`
}

type LoginResponse struct {
	models.User
	IsFirstLogin uint   `json:"isFirstLogin"`
	Token        string `json:"token"`
}

type LogoutForm struct {
	Platform *string `json:"platform" comment:"平台" binding:"required,oneof=twitta video_server"`
}

type UserUpdateForm struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Introduce string `json:"introduce"`
	FaceImg   string `json:"faceImg"`
}

type UserDetail struct {
	UserId      string `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Introduce   string `json:"introduce"`
	WechatCount int64  `json:"whatCount"`
	FansCount   int64  `json:"fansCount"`
	CreatedAt   string `json:"createdAt"`
}

type UserSearch struct {
	utils.PageList
	Records []*UserDetail `json:"records"`
}

type Face struct {
	models.User
	IsFirstLogin uint   `json:"isFirstLogin"`
	Token        string `json:"token"`
	IsFound      bool   `json:"isFound"`
}
