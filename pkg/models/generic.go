package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
)

// return collection, 丢弃interface设计，使用一份代码操作所有表
func GetCollection(db *mongo.Client, name string, collection string) *mongo.Collection {
	return db.Database(name).Collection(collection)
}

func GInsertOne[T any](ctx context.Context, collection *mongo.Collection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GInsertMany[T any](ctx context.Context, collection *mongo.Collection, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := collection.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GWhereFirst[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	var result T
	err := collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GWhereFind[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]*T, error) {
	var result []*T
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GWhereUpdate[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func GWhereDelete[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

func GWhereCount[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GPaginatorOrder[T any](ctx context.Context, collection *mongo.Collection, params *ListPageInput, sort string, filter interface{}) ([]*T, int64, int64, error) {
	var result []*T

	findOptions := new(options.FindOptions)
	if params != nil && params.Size > 0 && params.Page > 0 {
		findOptions.SetSkip(int64((params.Page - 1) * params.Size))
		findOptions.SetLimit(int64(params.Size))
	}

	if sort != "" {
		findOptions.SetSort(sort)
	}

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	if err = cur.All(ctx, &result); err != nil {
		return nil, 0, 0, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	pages := int64(math.Ceil(float64(count) / float64(params.Size)))
	return result, count, pages, nil
}
