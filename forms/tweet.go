package forms

import "twitta/pkg/utils"

type TweetCreateForm struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type Tweet struct {
	UserId        string   `json:"userId"`
	Username      string   `json:"username"`
	Avatar        string   `json:"avatar"`
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Images        []string `json:"images"`
	CreatedAt     string   `json:"createdAt"`
	ThumbCount    int64    `json:"thumbCount"`
	CommentCount  int64    `json:"commentCount"`
	RetweetCount  int64    `json:"retweetCount"`
	RetweetOf     string   `json:"retweetOf"`
	OriginTweet   *Tweet   `json:"originTweet,omitempty"`
}

type TweetList struct {
	utils.PageList
	Records []*Tweet `json:"records"`
}

type TweetFavoriteForm struct {
	Id string `json:"id" binding:"required"`
}
