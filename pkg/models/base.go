package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `bson:"created-at"`
	UpdatedAt time.Time `bson:"updated-at"`
	DeletedAt time.Time `bson:"deleted-at"`
}
