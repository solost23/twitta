package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Interface interface {
	TableName() string
	InsertOne(context.Context, *mongo.Client, string, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(context.Context, *mongo.Client, string, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	FindOne(context.Context, *mongo.Client, string, interface{}, interface{}, ...*options.FindOneOptions) error
	Find(context.Context, *mongo.Client, string, interface{}, interface{}, ...*options.FindOptions) error
	Update(context.Context, *mongo.Client, string, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Delete(context.Context, *mongo.Client, string, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
