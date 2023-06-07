package models

import (
	"time"
)

// 定义基本模型，对数据进行增删改查时结构可以变动
type User struct {
	BaseModel
	ID            string    `json:"id" bson:"_id"`
	Username      string    `json:"username" bson:"username"`
	Password      string    `json:"password" bson:"password"`
	Nickname      string    `json:"nickname" bson:"nickname"`
	Mobile        string    `json:"mobile" bson:"mobile"`
	Role          string    `json:"role" bson:"role"`
	Avatar        string    `json:"avatar" bson:"avatar"`
	Introduce     string    `json:"introduce" bson:"introduce"`
	Email         string    `json:"email" bson:"email"`
	FansCount     int64     `json:"fansCount" bson:"fans_count" comment:"关注数"`
	WechatCount   int64     `json:"wechatCount" bson:"wechat_count" comment:"关注数"`
	LastLoginTime time.Time `json:"lastLoginTime" bson:"last_login_time" comment:"用户最近登陆时间"`
	Disabled      uint      `json:"disabled" bson:"disabled" comment:"是否禁用用户 0: 非禁用 1: 禁用"`
	FaceImg       string    `json:"faceImg" bson:"face_img" comment:"用户人脸"`
	FaceEncoding  string    `json:"faceEncoding" bson:"face_encoding" comment:"用户人脸编码"`
}

func NewUser() *User {
	return &User{}
}

func (*User) TableName() string {
	return "users"
}

// func (t *User) InsertOne(ctx context.Context, db *mongo.Client, dbName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertOne(ctx, document, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *User) InsertMany(ctx context.Context, db *mongo.Client, dbName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertMany(ctx, documents, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *User) FindOne(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
//	err := db.Database(dbName).Collection(t.TableName()).FindOne(ctx, filter, opts...).Decode(result)
//	if err != nil {
//		return err
//	}
//	return nil
// }
//
// func (t *User) Find(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
//	cur, err := db.Database(dbName).Collection(t.TableName()).Find(ctx, filter, opts...)
//	if err != nil {
//		return err
//	}
//	if err = cur.All(ctx, results); err != nil {
//		return err
//	}
//	return nil
// }
//
// func (t *User) Update(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	updateResult, err := db.Database(dbName).Collection(t.TableName()).UpdateMany(ctx, filter, update, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return updateResult, nil
// }
//
// func (t *User) Delete(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	deleteResult, err := db.Database(dbName).Collection(t.TableName()).DeleteMany(ctx, filter, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return deleteResult, nil
// }
//
// func (t *User) Count(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
//	count, err := db.Database(dbName).Collection(t.TableName()).CountDocuments(ctx, filter, opts...)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
// }
