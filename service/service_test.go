package service

import (
	"context"
	"errors"
	"testing"

	"delivery-service/mocks"
	"delivery-service/storage/cache"
	"delivery-service/storage/mongodb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// get 1 active campaign from cache, campaign has no rules and validate data - cache hit
func TestGetCampaignsFromCache1(t *testing.T) {

	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			var campaigns []cache.Campaign
			campaigns = append(campaigns, cache.Campaign{Cid: "cid", Img: "image", Cta: "cta", IsActive: true, NoRestrictions: true, Rules: cache.Rule{}})
			return campaigns, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "c")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}

// get 1 active campaign from cache, campaign has rules and validate data - cache hit
func TestGetCampaignsFromCache2(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			var campaigns []cache.Campaign
			campaigns = append(campaigns, cache.Campaign{Cid: "cid", Img: "image", Cta: "cta", IsActive: true, NoRestrictions: true, Rules: cache.Rule{Os: map[string]bool{"android": true}, App: map[string]bool{"a": true}}})
			return campaigns, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "android")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}

// get 1 active campaign from cache and no data - cache hit
func TestGetCampaignsFromCache3(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			var campaigns []cache.Campaign
			return campaigns, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "ios")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(campaigns))
}

// get 0 active campaign from cache and campaign has rules and doesn't match requirements from request - cache hit
func TestGetCampaignsFromCache4(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			var campaigns []cache.Campaign
			campaigns = append(campaigns, cache.Campaign{Cid: "cid", Img: "image", Cta: "cta", IsActive: true, NoRestrictions: false, Rules: cache.Rule{Os: map[string]bool{"android": true, "ios": false}, App: map[string]bool{"a": true}}})
			return campaigns, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "ios")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(campaigns))
}

// get campaign from cache and campaign has rules and validate data - cache hit
func TestGetCampaignsFromCache5(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			var campaigns []cache.Campaign
			campaigns = append(campaigns, cache.Campaign{Cid: "cid", Img: "image", Cta: "cta", IsActive: true, NoRestrictions: false, Rules: cache.Rule{Os: map[string]bool{"android": true, "ios": true}, App: map[string]bool{"a": true}}})
			return campaigns, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "ios")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}

// get 1 active campaign from db, campaign has no rules and validate data - cache miss
func TestGetCampaignsFromDB1(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			return nil, errors.New("some error")
		},
	}
	mongodb.MongoInstance = mocks.MongoMock{
		FindMock: func(context.Context, string, interface{}, ...*options.FindOptions) (*mongo.Cursor, error) {
			cursor, _ := mongo.NewCursorFromDocuments(nil, nil, nil)
			return cursor, nil
		},
		AggregateMock: func(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			var data []interface{}
			data = append(data, mongodb.CampaignResponse{Id: "cid", Name: "name", Image: "image", Cta: "cta", IsActive: true, NoRestrictions: true, Rules: mongodb.RuleResponse{}})
			cursor, _ := mongo.NewCursorFromDocuments(data, nil, nil)
			return cursor, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "c")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}

// get campaigns from db, but got error from mongodb - cache miss
func TestGetCampaignsFromDB2(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			return nil, errors.New("some error")
		},
	}
	mongodb.MongoInstance = mocks.MongoMock{
		FindMock: func(context.Context, string, interface{}, ...*options.FindOptions) (*mongo.Cursor, error) {
			cursor, _ := mongo.NewCursorFromDocuments(nil, nil, nil)
			return cursor, nil
		},
		AggregateMock: func(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			return nil, errors.New("some error")
		},
	}
	svc := NewService()

	_, err := svc.GetCampaigns(context.Background(), "a", "b", "c")
	assert.Error(t, err)
}

// get campaigns from db, campaign has rules - cache miss
func TestGetCampaignsFromDB3(t *testing.T) {
	cache.CacheInstance = mocks.CacheMock{
		GetCampaignsByCountryMock: func(string) ([]cache.Campaign, error) {
			return nil, errors.New("some error")
		},
	}
	mongodb.MongoInstance = mocks.MongoMock{
		FindMock: func(context.Context, string, interface{}, ...*options.FindOptions) (*mongo.Cursor, error) {
			cursor, _ := mongo.NewCursorFromDocuments(nil, nil, nil)
			return cursor, nil
		},
		AggregateMock: func(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
			var data []interface{}
			data = append(data, mongodb.CampaignResponse{Id: "cid", Name: "name", Image: "image", Cta: "cta", IsActive: true, NoRestrictions: false, Rules: mongodb.RuleResponse{IncludeOs: []string{"android", "ios"}, IncludeApp: []string{"a"}}})
			cursor, _ := mongo.NewCursorFromDocuments(data, nil, nil)
			return cursor, nil
		},
	}
	svc := NewService()

	campaigns, err := svc.GetCampaigns(context.Background(), "a", "b", "ios")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, []Campaign{{Cid: "cid", Img: "image", Cta: "cta"}}, campaigns)
}
