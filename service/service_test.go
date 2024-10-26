package service

import (
	"context"
	"delivery-service/mocks"
	"delivery-service/storage/mongodb"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var iMongoDb mongodb.IMongoDb
var iMongoCollection mongodb.IMongoCollection

// get campaign from mongodb - success
func TestGetCampaigns1(t *testing.T) {

	mongodb.MongoDB = mocks.MongoMock{
		GetDbMock: func(db_name string) mongodb.IMongoDb {
			return iMongoDb
		},
	}
	iMongoDb = mocks.MongoDbMock{
		GetCollectionMock: func(coll_name string) mongodb.IMongoCollection {
			return iMongoCollection
		},
	}
	iMongoCollection = mocks.MongoCollectionMock{
		FindOneMock: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
			data := bson.M{
				"rules": bson.A{"app", "country", "os"},
			}
			result := mongo.NewSingleResultFromDocument(data, nil, nil)
			return result, nil
		},
		AggregateMock: func(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			data := bson.A{
				bson.M{
					"_id":   "cid",
					"image": "image",
					"cta":   "cta",
				},
			}
			cursor, _ := mongo.NewCursorFromDocuments(data, nil, nil)
			return cursor, nil
		},
	}

	svc := NewService()

	params := map[string]string{"app": "a", "country": "b", "os": "c"}

	campaigns, err := svc.GetCampaigns(context.Background(), params, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}

// get campaign from mongodb - failed because mongo return some error
func TestGetCampaigns2(t *testing.T) {
	mongodb.MongoDB = mocks.MongoMock{
		GetDbMock: func(db_name string) mongodb.IMongoDb {
			return iMongoDb
		},
	}
	iMongoDb = mocks.MongoDbMock{
		GetCollectionMock: func(coll_name string) mongodb.IMongoCollection {
			return iMongoCollection
		},
	}
	iMongoCollection = mocks.MongoCollectionMock{
		FindOneMock: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
			return nil, errors.New("some error")
		},
	}

	svc := NewService()

	params := map[string]string{"app": "a", "country": "b", "os": "c"}

	_, err := svc.GetCampaigns(context.Background(), params, 10, 0)
	assert.Error(t, err)
}

// get campaign from mongodb - failed because mongo return some error
func TestGetCampaigns3(t *testing.T) {
	mongodb.MongoDB = mocks.MongoMock{
		GetDbMock: func(db_name string) mongodb.IMongoDb {
			return iMongoDb
		},
	}
	iMongoDb = mocks.MongoDbMock{
		GetCollectionMock: func(coll_name string) mongodb.IMongoCollection {
			return iMongoCollection
		},
	}
	iMongoCollection = mocks.MongoCollectionMock{
		FindOneMock: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
			data := bson.M{
				"rules": bson.A{"app", "country", "os"},
			}
			result := mongo.NewSingleResultFromDocument(data, nil, nil)
			return result, nil
		},
		AggregateMock: func(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			return nil, errors.New("some error")
		},
	}

	svc := NewService()

	params := map[string]string{"app": "a", "country": "b", "os": "c"}

	_, err := svc.GetCampaigns(context.Background(), params, 10, 0)
	assert.Error(t, err)
}

// get campaign from mongodb - failed because unknown rules parameters passed
func TestGetCampaigns4(t *testing.T) {
	mongodb.MongoDB = mocks.MongoMock{
		GetDbMock: func(db_name string) mongodb.IMongoDb {
			return iMongoDb
		},
	}
	iMongoDb = mocks.MongoDbMock{
		GetCollectionMock: func(coll_name string) mongodb.IMongoCollection {
			return iMongoCollection
		},
	}
	iMongoCollection = mocks.MongoCollectionMock{
		FindOneMock: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
			data := bson.M{
				"rules": bson.A{"app", "country", "os"},
			}
			result := mongo.NewSingleResultFromDocument(data, nil, nil)
			return result, nil
		},
	}

	svc := NewService()

	params := map[string]string{"app": "a", "country": "b", "os": "c", "unknown": "unknown"}

	_, err := svc.GetCampaigns(context.Background(), params, 10, 0)
	assert.Error(t, err)
}

// get campaign from mongodb - success, accept new rule parameter state
func TestGetCampaigns5(t *testing.T) {
	mongodb.MongoDB = mocks.MongoMock{
		GetDbMock: func(db_name string) mongodb.IMongoDb {
			return iMongoDb
		},
	}
	iMongoDb = mocks.MongoDbMock{
		GetCollectionMock: func(coll_name string) mongodb.IMongoCollection {
			return iMongoCollection
		},
	}
	iMongoCollection = mocks.MongoCollectionMock{
		FindOneMock: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
			data := bson.M{
				"rules": bson.A{"app", "country", "os", "state"},
			}
			result := mongo.NewSingleResultFromDocument(data, nil, nil)
			return result, nil
		},
		AggregateMock: func(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			data := bson.A{
				bson.M{
					"_id":   "cid",
					"image": "image",
					"cta":   "cta",
				},
			}
			cursor, _ := mongo.NewCursorFromDocuments(data, nil, nil)
			return cursor, nil
		},
	}

	svc := NewService()

	params := map[string]string{"app": "a", "country": "b", "os": "c", "state": "d"}

	campaigns, err := svc.GetCampaigns(context.Background(), params, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}
