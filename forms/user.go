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
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Device   string `json:"device" comment:"设备类型" binding:"required,oneof=ios android web"`
}

type LoginResponse struct {
	models.User
	IsFirstLogin uint   `json:"isFirstLogin"`
	Token        string `json:"token"`
}

type LogoutForm struct {
	Device string `json:"device" comment:"设备类型" binding:"required,oneof=ios android web"`
}

type UserUpdateForm struct {
	Username  string `json:"username" binding:"required"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Introduce string `json:"introduce"`
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
