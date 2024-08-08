package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"math"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 丢弃interface设计，使用一份代码操作所有表

func GInsertOne[T TableNamer](ctx context.Context, db *mongo.Database, document T) error {
	var t T
	collection := db.Collection(t.TableName())

	_, err := collection.InsertOne(ctx, document)
	return err
}

func GInsertMany[T TableNamer](ctx context.Context, db *mongo.Database, documents []T) error {
	var t T
	collection := db.Collection(t.TableName())

	anySlice := make([]any, 0, len(documents))
	for i := 0; i != len(documents); i++ {
		anySlice = append(anySlice, documents[i])
	}
	_, err := collection.InsertMany(ctx, anySlice)
	return err
}

func GWhereFirst[T TableNamer](ctx context.Context, db *mongo.Database, filter bson.M) (T, error) {
	var result T

	var t T
	collection := db.Collection(t.TableName())

	err := collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func GWhereFind[T TableNamer](ctx context.Context, db *mongo.Database, filter bson.M) ([]T, error) {
	var result []T

	var t T
	collection := db.Collection(t.TableName())

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GWhereUpdate[T TableNamer](ctx context.Context, db *mongo.Database, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	var t T
	collection := db.Collection(t.TableName())

	return collection.UpdateMany(ctx, filter, update)
}

func GWhereDelete[T TableNamer](ctx context.Context, db *mongo.Database, filter bson.M) (*mongo.DeleteResult, error) {
	var t T
	collection := db.Collection(t.TableName())

	result, err := collection.DeleteMany(ctx, filter)
	return result, err
}

func GWhereCount[T TableNamer](ctx context.Context, db *mongo.Database, filter bson.M) (int64, error) {
	var t T
	collection := db.Collection(t.TableName())

	total, err := collection.CountDocuments(ctx, filter)
	return total, err
}

func GPaginatorOrder[T TableNamer](ctx context.Context, db *mongo.Database, params *ListPageInput, sort bson.M, filter bson.M) ([]T, int64, int64, error) {
	page := params.Page
	size := params.Size

	var result []T

	opts := new(options.FindOptions)
	opts.SetSkip(int64((page - 1) * size))
	opts.SetLimit(int64(size))
	opts.SetSort(sort)

	var t T
	collection := db.Collection(t.TableName())
	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, 0, err
	}
	if err = cur.All(ctx, &result); err != nil {
		return nil, 0, 0, err
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	return result, total, int64(math.Ceil(float64(total) / float64(params.Size))), nil
}
