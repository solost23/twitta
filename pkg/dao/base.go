package dao

type ListPageInput struct {
	Size int `comment:"页长"`
	Page int `comment:"当前页"`
}

type ListPageOutput struct {
	Size    int   `json:"size"`
	Pages   int64 `json:"pages"`
	Total   int64 `json:"total"`
	Current int   `json:"current"`
}
