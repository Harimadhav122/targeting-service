package service

import (
	"context"
	"delivery-service/metrics"
	"delivery-service/storage"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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
	cache storage.ICampaignCache
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
		cache: storage.NewCache(),
	}
}

// GetCampaigns implements the business logic
func (s *campaignService) GetCampaigns(ctx context.Context, app, country, os string) ([]Campaign, error) {

	start := time.Now()
	var result []Campaign
	var campaigns []storage.Campaign
	var err error

	//get campaigns from cache
	campaigns, err = s.cache.GetCampaignsByCountry(country)
	if err != nil {
		// cache miss
		metrics.CacheMiss.Add(1)
		level.Warn(logger).Log("method", "GetCampaigns", "err", err, "took", time.Since(start), "cache", "miss")
		// get the campaigns from db
		return nil, err
	}
	// cache hit
	metrics.CacheHit.Add(1)
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

	level.Info(logger).Log("method", "GetCampaigns", "took", time.Since(start))
	return result, nil
}
