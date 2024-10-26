package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"delivery-service/endpoints"
	"delivery-service/mocks"
	"delivery-service/service"
	"delivery-service/storage/mongodb"
	"delivery-service/transport"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var iMongoDb mongodb.IMongoDb
var iMongoCollection mongodb.IMongoCollection

// test 200 http status code
func TestMain1(t *testing.T) {

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

	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us&os=android&limit=10&page=0")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// test 400 http status code by missing os param
func TestMain2(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing country param
func TestMain3(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&os=android")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing app param
func TestMain4(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?country=us&os=ios")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing limit param
func TestMain5(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us&os=android")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing page param
func TestMain6(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us&os=android&limit=10")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by passing unknown param
func TestMain7(t *testing.T) {
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

	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us&os=android&unknown=unknown&limit=10&page=0")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 405 http status code by unsupported method
func TestMain8(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	data := "{}"
	reader := bytes.NewReader([]byte(data))

	// Make a valid request to the test server
	resp, err := http.Post(server.URL+"/v1/delivery?app=com.gametion.ludokinggame&country=us", "application/json", reader)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
