package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"createdAt" bson:"created-at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated-at"`
	DeletedAt time.Time `json:"deletedAt" bson:"deleted-at"`
}

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
