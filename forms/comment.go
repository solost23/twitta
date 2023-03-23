package forms

import "twitta/pkg/utils"

type CommentInsertForm struct {
	Content  string `json:"content" binding:"required,min=3"`
	ParentId string `json:"parentId"`
}

type CommentListResponse struct {
	Id        string `json:"id"`
	PID       string `json:"pid"`
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Introduce string `json:"introduce"`
	Content   string `json:"content"`
	Children  []*CommentListResponse
}

func (c *CommentListResponse) GetPID() string {
	return c.PID
}

func (c *CommentListResponse) GetID() string {
	return c.Id
}

func (c *CommentListResponse) AppendChildrenNode(node utils.TreeNode) {
	c.Children = append(c.Children, node.(*CommentListResponse))
}
