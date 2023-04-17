package forms

import "twitta/pkg/utils"

type CommentInsertForm struct {
	Content  *string `json:"content" binding:"required,min=3"`
	ParentId *string `json:"parentId"`
	utils.PageForm
}

type Comment struct {
	Id        *string `json:"id"`
	PID       *string `json:"pid"`
	UserId    *string `json:"userId"`
	Username  *string `json:"username"`
	Avatar    *string `json:"avatar"`
	Introduce *string `json:"introduce"`
	Content   *string `json:"content"`
	Children  []*Comment
}

type CommentList struct {
	Record []*Comment `json:"records"`
	*utils.PageList
}

func (c *Comment) GetPID() string {
	return *c.PID
}

func (c *Comment) GetID() string {
	return *c.Id
}

func (c *Comment) AppendChildrenNode(node utils.TreeNode) {
	c.Children = append(c.Children, node.(*Comment))
}
