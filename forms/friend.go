package forms

type FriendApplicationSendForm struct {
	UserId  string `json:"userId" binding:"required"`
	Content string `json:"content" binding:"required,min=2"`
}

type FriendApplicationListResponse struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Type      uint   `json:"type"`
	CreatedAt string `json:"createdAt"`
}
