package mocks

import (
	"context"
	"delivery-service/storage/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMock struct {
	GetDbMock func(string) mongodb.IMongoDb
}

type MongoDbMock struct {
	GetCollectionMock func(string) mongodb.IMongoCollection
}

type MongoCollectionMock struct {
	FindOneMock   func(context.Context, interface{}, ...*options.FindOneOptions) (*mongo.SingleResult, error)
	AggregateMock func(context.Context, interface{}, ...*options.AggregateOptions) (*mongo.Cursor, error)
}

func (m MongoMock) GetDb(db_name string) mongodb.IMongoDb {
	return m.GetDbMock(db_name)
}

func (m MongoDbMock) GetCollection(coll_name string) mongodb.IMongoCollection {
	return m.GetCollectionMock(coll_name)
}

func (m MongoCollectionMock) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	return m.FindOneMock(ctx, filter, opts...)
}

func (m MongoCollectionMock) Aggregate(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.AggregateMock(ctx, filter, opts...)
}
