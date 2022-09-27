package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 定义基本模型，对数据进行增删改查时结构可以变动
type User struct {
	BaseModel
	ID            string    `bson:"_id"`
	Username      string    `bson:"username"`
	Password      string    `bson:"password"`
	Nickname      string    `bson:"nickname"`
	Mobile        string    `bson:"mobile"`
	Role          string    `bson:"role"`
	Avatar        string    `bson:"avatar"`
	Introduce     string    `bson:"introduce"`
	FansCount     int64     `bson:"fans-count"`
	LastLoginTime time.Time `bson:"last-login-time" comment:"用户最近登陆时间"`
	Disabled      uint      `bson:"disabled" comment:"是否禁用用户 0: 非禁用 1: 禁用"`
}

func NewUser() Interface {
	return &User{}
}

func (*User) TableName() string {
	return "users"
}

func (*User) InsertOne(ctx context.Context, db *mongo.Client, dbName string, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := db.Database(dbName).Collection(collection).InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*User) InsertMany(ctx context.Context, db *mongo.Client, dbName string, collection string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := db.Database(dbName).Collection(collection).InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*User) FindOne(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	err := db.Database(dbName).Collection(collection).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (*User) Find(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	cur, err := db.Database(dbName).Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (*User) Update(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := db.Database(dbName).Collection(collection).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func (*User) Delete(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := db.Database(dbName).Collection(collection).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}
