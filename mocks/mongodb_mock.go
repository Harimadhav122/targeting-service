package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMock struct {
	FindMock      func(context.Context, string, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	AggregateMock func(context.Context, string, interface{}, ...*options.AggregateOptions) (*mongo.Cursor, error)
}

func (m MongoMock) Find(ctx context.Context, coll string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.FindMock(ctx, coll, filter, opts...)
}

func (m MongoMock) Aggregate(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.AggregateMock(ctx, collection, filter, opts...)
}
