package forms

type FriendApplicationSendForm struct {
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type FriendApplicationListResponse struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Type      uint   `json:"type"`
	CreatedAt string `json:"createdAt"`
}
