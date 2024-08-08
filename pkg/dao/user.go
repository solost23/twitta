package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updated_at,omitempty"`
	DeletedAt     time.Time          `json:"deletedAt" bson:"deleted_at,omitempty"`
	Username      string             `json:"username" bson:"username,omitempty" comment:"用户名"`
	Password      string             `json:"password" bson:"password,omitempty" comment:"密码"`
	Nickname      string             `json:"nickname" bson:"nickname,omitempty" comment:"昵称"`
	Mobile        string             `json:"mobile" bson:"mobile,omitempty" comment:"手机"`
	Role          string             `json:"role" bson:"role,omitempty" comment:"角色"`
	Avatar        string             `json:"avatar" bson:"avatar,omitempty" comment:"头像"`
	Introduce     string             `json:"introduce" bson:"introduce,omitempty" comment:"介绍"`
	Email         string             `json:"email" bson:"email,omitempty" comment:"邮箱"`
	FansCount     int64              `json:"fansCount" bson:"fans_count,omitempty" comment:"关注数"`
	WechatCount   int64              `json:"wechatCount" bson:"wechat_count,omitempty" comment:"关注数"`
	LastLoginTime time.Time          `json:"lastLoginTime" bson:"last_login_time,omitempty" comment:"用户最近登陆时间"`
	Disabled      uint               `json:"disabled" bson:"disabled,omitempty" comment:"是否禁用用户 0: 非禁用 1: 禁用"`
	FaceImg       string             `json:"faceImg" bson:"face_img,omitempty" comment:"用户人脸"`
	FaceEncoding  string             `json:"faceEncoding" bson:"face_encoding,omitempty" comment:"用户人脸编码"`
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

//func (m *User) Conn() *mongo.Collection {
//	return NewCollection(m.TableName()).Build()
//}
