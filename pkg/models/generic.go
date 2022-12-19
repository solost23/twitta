package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 丢弃interface设计，使用一份代码操作所有表
func GInsertOne(ctx context.Context, db *mongo.Client, dbName string, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := db.Database(dbName).Collection(collection).InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GInsertMany(ctx context.Context, db *mongo.Client, dbName string, collection string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := db.Database(dbName).Collection(collection).InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GFindOne(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	err := db.Database(dbName).Collection(collection).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func GFind(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	cur, err := db.Database(dbName).Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func GUpdate(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := db.Database(dbName).Collection(collection).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func GDelete(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := db.Database(dbName).Collection(collection).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

func GCount(ctx context.Context, db *mongo.Client, dbName string, collection string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := db.Database(dbName).Collection(collection).CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}
