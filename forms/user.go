package forms

import "Twitta/pkg/models"

type RegisterForm struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Nickname  string `json:"nickname"`
	Mobile    string `json:"mobile"`
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
