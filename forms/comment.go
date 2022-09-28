package forms

type CommentInsertForm struct {
	Content  string `json:"content" binding:"required,min=3"`
	ParentId string `json:"parentId"`
}

type CommentListResponse struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Introduce string `json:"introduce"`
	Id        string `json:"id"`
	Content   string `json:"content"`
}
