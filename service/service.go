package service

import (
	"context"
	"delivery-service/metrics"
	"delivery-service/utils"
	"os"
	"time"

	"delivery-service/storage/cache"
	"delivery-service/storage/mongodb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
)

// Campaign represents a campaign entity
type Campaign struct {
	Cid string `json:"cid"`
	Img string `json:"img"`
	Cta string `json:"cta"`
}

// Service defines the behavior of our campaign service
type Service interface {
	GetCampaigns(ctx context.Context, app string, country string, os string) ([]Campaign, error)
}

// campaignService is the implementation of the Service interface
type campaignService struct {
	cache cache.ICampaignCache
	mongo mongodb.IMongo
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
		cache: cache.NewCache(),
		mongo: mongodb.NewMongo(),
	}
}

// GetCampaigns implements the business logic
func (s *campaignService) GetCampaigns(ctx context.Context, app, country, os string) ([]Campaign, error) {

	start := time.Now()
	var campaigns []cache.Campaign
	var err error

	//get campaigns from cache
	campaigns, err = s.cache.GetCampaignsByCountry(country)
	if err != nil {
		// cache miss
		metrics.CacheMiss.Add(1)

		// get the campaigns from db
		filter := getFilterForCampaignsPerCountry(country)
		cursor, err := s.mongo.Aggregate(ctx, "campaigns_details", filter)

		var camps []mongodb.CampaignResponse
		if err = cursor.All(context.TODO(), &camps); err != nil {
			level.Error(logger).Log("method", "GetCampaigns", err, "error decoding cursor")
			return nil, err
		}

		result := verifyCampaignsFromDB(camps, app, os)

		level.Info(logger).Log("method", "GetCampaigns", "took", time.Since(start))
		return result, nil
	}
	// cache hit
	metrics.CacheHit.Add(1)
	result := verifyCampaignsFromCache(campaigns, app, os)

	level.Info(logger).Log("method", "GetCampaigns", "took", time.Since(start))
	return result, nil
}

func verifyCampaignsFromCache(campaigns []cache.Campaign, app, os string) []Campaign {
	var result []Campaign

	for _, campaign := range campaigns {
		if !campaign.IsActive {
			continue
		}
		if campaign.NoRestrictions {
			result = append(result, Campaign{campaign.Cid, campaign.Img, campaign.Cta})
			continue
		} else {
			if value, ok := campaign.Rules.Os[os]; ok && !value {
				continue
			}
			if value, ok := campaign.Rules.App[app]; ok && !value {
				continue
			}
			result = append(result, Campaign{campaign.Cid, campaign.Img, campaign.Cta})
		}
	}
	return result
}

func verifyCampaignsFromDB(campaigns []mongodb.CampaignResponse, app, os string) []Campaign {

	var result []Campaign

	for _, campaign := range campaigns {
		if campaign.NoRestrictions {
			result = append(result, Campaign{campaign.Id, campaign.Image, campaign.Cta})
			continue
		} else {
			includeOsLen := len(campaign.Rules.IncludeOs)
			excludeOsLen := len(campaign.Rules.ExcludeOs)

			includeAppLen := len(campaign.Rules.IncludeApp)
			excludeAppLen := len(campaign.Rules.ExcludeApp)

			if includeOsLen > 0 && !utils.Contains(campaign.Rules.IncludeOs, os) {
				continue
			}
			if excludeOsLen > 0 && utils.Contains(campaign.Rules.ExcludeOs, os) {
				continue
			}
			if includeAppLen > 0 && !utils.Contains(campaign.Rules.IncludeApp, app) {
				continue
			}
			if excludeAppLen > 0 && utils.Contains(campaign.Rules.ExcludeApp, app) {
				continue
			}

			result = append(result, Campaign{campaign.Id, campaign.Image, campaign.Cta})
		}
	}
	return result
}

func getFilterForCampaignsPerCountry(country string) bson.A {

	var pipeline bson.A

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"isActive": true,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"$and": bson.A{
				bson.M{
					"$or": bson.A{
						bson.M{
							"rules.includeCountry": bson.M{
								"$size": 0,
							},
						},
						bson.M{
							"$and": bson.A{
								bson.M{
									"rules.includeCountry": bson.M{
										"$in": bson.A{country},
									},
								},
							},
						},
					},
				},
				bson.M{
					"$or": bson.A{
						bson.M{
							"rules.excludeCountry": bson.M{
								"$size": 0,
							},
						},
						bson.M{
							"$and": bson.A{
								bson.M{
									"rules.excludeCountry": bson.M{
										"$not": bson.M{
											"$in": bson.A{country},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	return pipeline
}
