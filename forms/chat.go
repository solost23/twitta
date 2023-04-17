package forms

import (
	"twitta/pkg/utils"
)

type ChatList struct {
	*utils.PageList
	Records []*Chat `json:"records"`
}

type Chat struct {
	Msg       *string `json:"msg"`
	UserId    *string `json:"userId"`
	CreatedAt *string `json:"createdAt"`
}
