package forms

import "twitta/pkg/utils"

type TweetCreateForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Tweet struct {
	// 发推人ID 名称 头像 推文ID 标题，内容，发推时间 点赞数 评论数
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	ID           string `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CreatedAt    string `json:"createdAt"`
	ThumbCount   int64  `json:"thumbCount"`
	CommentCount int64  `json:"commentCount"`
}

type TweetList struct {
	utils.PageList
	Records []*Tweet `json:"records"`
}

type TweetFavoriteForm struct {
	Id string `json:"id" binding:"required"`
}
