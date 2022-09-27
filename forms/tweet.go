package forms

type TweetCreateForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TweetListResponse struct {
	// 发推人ID 名称 头像 推文ID 标题，内容，发推时间 点赞数 评论数
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	TweetId      string `json:"tweetId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	TweetTime    string `json:"tweetTime"`
	ThumbCount   int64  `json:"thumbCount"`
	CommentCount int64  `json:"commentCount"`
}

type TweetFavoriteForm struct {
	Id string `json:"id"`
}
