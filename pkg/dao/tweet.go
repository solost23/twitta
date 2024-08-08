package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tweet struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updated_at,omitempty"`
	DeletedAt    time.Time          `json:"deletedAt" bson:"deleted_at,omitempty"`
	UserID       string             `json:"userId" bson:"user_id,omitempty" comment:"用户 ID"`
	Title        string             `json:"title" bson:"title,omitempty" comment:"标题"`
	Content      string             `json:"content" bson:"content,omitempty" comment:"内容"`
	ThumbCount   int64              `json:"thumbCount" bson:"thumb_count,omitempty" comment:"点赞数"`
	CommentCount int64              `json:"commentCount" bson:"comment_count,omitempty" comment:"评论数"`
}

//	func NewTweet() *Tweet {
//		return &Tweet{}
//	}
func (*Tweet) TableName() string {
	return "tweets"
}

// func (t *Tweet) InsertOne(ctx context.Context, db *mongo.Client, dbName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertOne(ctx, document, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *Tweet) InsertMany(ctx context.Context, db *mongo.Client, dbName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertMany(ctx, documents, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *Tweet) FindOne(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
//	err := db.Database(dbName).Collection(t.TableName()).FindOne(ctx, filter, opts...).Decode(result)
//	if err != nil {
//		return err
//	}
//	return nil
// }
//
// func (t *Tweet) Find(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
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
// func (t *Tweet) Update(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	updateResult, err := db.Database(dbName).Collection(t.TableName()).UpdateMany(ctx, filter, update, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return updateResult, nil
// }
//
// func (t *Tweet) Delete(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	deleteResult, err := db.Database(dbName).Collection(t.TableName()).DeleteMany(ctx, filter, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return deleteResult, nil
// }
//
// func (t *Tweet) Count(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
//	count, err := db.Database(dbName).Collection(t.TableName()).CountDocuments(ctx, filter, opts...)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
// }

//func (m *Tweet) Conn() *mongo.Collection {
//	return NewCollection(m.TableName()).Build()
//}
