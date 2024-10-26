package service

import (
	"context"
	"os"

	local_error "delivery-service/errors"
	"delivery-service/storage/mongodb"
	"delivery-service/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
)

// Campaign represents a campaign entity
type Campaign struct {
	Cid string `json:"cid" bson:"_id"`
	Img string `json:"img" bson:"image"`
	Cta string `json:"cta" bson:"cta"`
}

type Parameters struct {
	Rules []string `bson:"rules"`
}

// Service defines the behavior of our campaign service
type Service interface {
	GetCampaigns(ctx context.Context, params map[string]string, limit, offset int) ([]Campaign, error)
}

// campaignService is the implementation of the Service interface
type campaignService struct {
	db mongodb.IMongoDb
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "service")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

}

// NewService creates and returns a new Campaign Service
func NewService() Service {
	return &campaignService{
		db: mongodb.MongoDB.GetDb("campaigns"),
	}
}

// GetCampaigns implements the business logic
func (s *campaignService) GetCampaigns(ctx context.Context, params map[string]string, limit, offset int) ([]Campaign, error) {

	var campaigns []Campaign
	var parameters Parameters
	var coll mongodb.IMongoCollection
	var err error

	coll = s.db.GetCollection("rules_parameters")
	result, err := coll.FindOne(ctx, bson.M{"_id": "current"})

	if err != nil {
		level.Error(logger).Log("method", "GetCampaigns", "msg", "mongodb findOne failed", "err", err)
		return nil, err
	}

	if err = result.Decode(&parameters); err != nil {
		level.Error(logger).Log("method", "GetCampaigns", "msg", "error decoding doc", "err", err)
		return nil, err
	}

	for param := range params {
		if !utils.Contains(parameters.Rules, param) {
			return nil, &local_error.ErrUnknownParams{Param: param}
		}
	}

	coll = s.db.GetCollection(params["country"])
	filter := getCampaignsFilter(params, limit, offset)
	cursor, err := coll.Aggregate(ctx, filter)

	if err != nil {
		level.Error(logger).Log("method", "GetCampaigns", "msg", "mongodb aggregate failed", "err", err)
		return nil, err
	}

	if err = cursor.All(context.TODO(), &campaigns); err != nil {
		level.Error(logger).Log("method", "GetCampaigns", err, "error decoding cursor")
		return nil, err
	}

	return campaigns, nil
}

func getCampaignsFilter(parameters map[string]string, limit, offset int) bson.A {

	var pipeline bson.A

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"_id": 1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "campaigns_details",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "result",
		},
	})

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"result.isActive": true,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$unwind": bson.M{
			"path":                       "$result",
			"preserveNullAndEmptyArrays": false,
		},
	})

	for param, paramValue := range parameters {
		includeParam := "result.rules.include" + param
		excludeParam := "result.rules.exclude" + param

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{
						"$or": bson.A{
							bson.M{
								includeParam: nil,
							},
							bson.M{
								includeParam: bson.M{
									"$in": bson.A{paramValue},
								},
							},
						},
					},
					bson.M{
						"$or": bson.A{
							bson.M{
								excludeParam: nil,
							},
							bson.M{
								excludeParam: bson.M{
									"$not": bson.M{
										"$in": bson.A{paramValue},
									},
								},
							},
						},
					},
				},
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$limit": limit,
	})

	pipeline = append(pipeline, bson.M{
		"$skip": limit * offset,
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"image": "$result.image",
			"cta":   "$result.cta",
		},
	})

	return pipeline
}
