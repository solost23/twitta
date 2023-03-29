package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"createdAt" bson:"created-at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated-at"`
	DeletedAt time.Time `json:"deletedAt" bson:"deleted-at"`
}
