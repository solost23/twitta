package forms

type ChatListResponse struct {
	Msg       string `json:"msg"`
	UserId    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}
