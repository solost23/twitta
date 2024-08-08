package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CommentTypeThumb = iota
	CommentTypeComment
)

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at,omitempty"`
	DeletedAt time.Time          `json:"deletedAt" bson:"deleted_at,omitempty"`
	UserId    string             `json:"userId" bson:"user_id,omitempty" comment:"用户 ID"`
	TweetId   string             `json:"tweetId" bson:"tweet_id,omitempty" comment:"推文 ID"`
	Content   string             `json:"content" bson:"content,omitempty" comment:"推文内容"`
	ParentId  string             `json:"parent" bson:"parent,omitempty" comment:"父节点 ID"`
	Type      uint               `json:"type" bson:"type,omitempty" comment:"评论类型 0: 点赞 1: 评论"`
}

func (*Comment) TableName() string {
	return "fans"
}

// func (t *Comment) InsertOne(ctx context.Context, db *mongo.Client, dbName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertOne(ctx, document, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *Comment) InsertMany(ctx context.Context, db *mongo.Client, dbName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertMany(ctx, documents, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
// }
//
// func (t *Comment) FindOne(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
//	err := db.Database(dbName).Collection(t.TableName()).FindOne(ctx, filter, opts...).Decode(result)
//	if err != nil {
//		return err
//	}
//	return nil
// }
//
// func (t *Comment) Find(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
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
// func (t *Comment) Update(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	updateResult, err := db.Database(dbName).Collection(t.TableName()).UpdateMany(ctx, filter, update, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return updateResult, nil
// }
//
// func (t *Comment) Delete(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	deleteResult, err := db.Database(dbName).Collection(t.TableName()).DeleteMany(ctx, filter, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return deleteResult, nil
// }
//
// func (t *Comment) Count(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
//	count, err := db.Database(dbName).Collection(t.TableName()).CountDocuments(ctx, filter, opts...)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
// }
