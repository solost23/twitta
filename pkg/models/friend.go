package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Friend struct {
	BaseModel
	ID       string `json:"id" bson:"_id"`
	UserId   string `json:"userId" bson:"user_id"`
	FriendId string `json:"friendId" bson:"friend_id"`
}

func NewFriend() Interface {
	return &Friend{}
}

func (*Friend) TableName() string {
	return "friends"
}

func (t *Friend) InsertOne(ctx context.Context, db *mongo.Client, dbName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := db.Database(dbName).Collection(t.TableName()).InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Friend) InsertMany(ctx context.Context, db *mongo.Client, dbName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := db.Database(dbName).Collection(t.TableName()).InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Friend) FindOne(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	err := db.Database(dbName).Collection(t.TableName()).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (t *Friend) Find(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	cur, err := db.Database(dbName).Collection(t.TableName()).Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (t *Friend) Update(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := db.Database(dbName).Collection(t.TableName()).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func (t *Friend) Delete(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := db.Database(dbName).Collection(t.TableName()).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

func (t *Friend) Count(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := db.Database(dbName).Collection(t.TableName()).CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}
