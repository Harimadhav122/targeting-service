package cache

import (
	"context"
	"os"
	"sync"

	"delivery-service/storage/mongodb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	supported_countries = [4]string{"us", "india", "canada", "germany"}
	supported_os        = [2]string{"android", "ios"}
	supported_apps      = [5]string{"com.abc.xyz", "com.gametion.ludokinggame", "com.apple.in", "com.google.in", "com.samsung.in"}
	ctx                 = context.Background()
	mongo               = mongodb.NewMongo()
	logger              log.Logger
	wg                  sync.WaitGroup
)

type ICampaignCache interface {
	GetCampaignsByCountry(country string) ([]Campaign, error)
}

type CampaignCache struct {
	mutex            sync.RWMutex
	Campaigns        map[string][]string // country -> cid
	CampaignsDetails map[string]Campaign // cid -> campaign_details
}

type Campaign struct {
	Cid            string
	Img            string
	Cta            string
	IsActive       bool
	NoRestrictions bool
	Rules          Rule
}

type Rule struct {
	Os  map[string]bool
	App map[string]bool
}

var campaignCache = &CampaignCache{Campaigns: make(map[string][]string), CampaignsDetails: make(map[string]Campaign)}
var CacheInstance ICampaignCache

func init() {

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "cache")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

	level.Info(logger).Log("method", "getAndFillCache", "msg", "cache fill started")

	// increment goroutine counter
	wg.Add(1)
	// fill cache with all campaigns details
	go getCampaignsAndFillCache()

	// fill cache with campaigns per country
	for _, country := range supported_countries {
		wg.Add(1)
		go getCampaignsByCountryAndFillCache(country)
	}
	// wait until all goroutines have finished
	wg.Wait()
	CacheInstance = campaignCache
	level.Info(logger).Log("method", "getAndFillCache", "msg", "cache fill completed")
}

func getCampaignsAndFillCache() {
	// goroutine execution completed
	defer wg.Done()
	var data []mongodb.CampaignResponse
	cursor, err := mongo.Find(ctx, "campaigns_details", bson.D{})
	if err != nil {
		level.Error(logger).Log("method", "getAndFillCache", "err", err)
		panic(err)
	}

	if err = cursor.All(context.TODO(), &data); err != nil {
		level.Error(logger).Log("method", "getAndFillCache", err, "error decoding cursor")
		panic(err)
	}
	for _, d := range data {
		rules := Rule{Os: make(map[string]bool), App: make(map[string]bool)}

		if len(d.Rules.IncludeOs) > 0 {
			for _, os := range d.Rules.IncludeOs {
				rules.Os[os] = true
			}
			for _, os := range supported_os {
				if _, ok := rules.Os[os]; !ok {
					rules.App[os] = false
				}
			}
		} else if len(d.Rules.ExcludeOs) > 0 {
			for _, os := range d.Rules.ExcludeOs {
				rules.Os[os] = false
			}
		} else {
			for _, os := range supported_os {
				rules.Os[os] = true
			}
		}

		if len(d.Rules.IncludeApp) > 0 {
			for _, app := range d.Rules.IncludeApp {
				rules.App[app] = true
			}
			for _, app := range supported_apps {
				if _, ok := rules.App[app]; !ok {
					rules.App[app] = false
				}
			}
		} else if len(d.Rules.ExcludeApp) > 0 {
			for _, app := range d.Rules.ExcludeApp {
				rules.App[app] = false
			}
		} else {
			for _, app := range supported_apps {
				rules.App[app] = true
			}
		}

		campaign := Campaign{Cid: d.Id, Img: d.Image, Cta: d.Cta, IsActive: d.IsActive, NoRestrictions: d.NoRestrictions, Rules: rules}
		setCampaignDetails(d.Id, campaign)
	}
}

func getCampaignsByCountryAndFillCache(country string) {
	// goroutine execution completed
	defer wg.Done()
	cursor, err := mongo.Find(ctx, country, bson.D{})
	if err != nil {
		level.Error(logger).Log("method", "getAndFillCache", "err", err)
		panic(err)
	}
	var result []string
	var data []mongodb.CampaignIdResponse

	if err = cursor.All(context.TODO(), &data); err != nil {
		level.Error(logger).Log("method", "getAndFillCache", err, "error decoding cursor")
		panic(err)
	}
	for _, d := range data {
		result = append(result, d.Id)
	}
	setCampaign(country, result)
}

func NewCache() ICampaignCache {
	return CacheInstance
}

func setCampaign(key string, value []string) {
	campaignCache.mutex.Lock()
	defer campaignCache.mutex.Unlock()
	campaignCache.Campaigns[key] = value
}

func getCampaign(key string) []string {
	campaignCache.mutex.RLock()
	defer campaignCache.mutex.RUnlock()
	return campaignCache.Campaigns[key]
}

func setCampaignDetails(key string, value Campaign) {
	campaignCache.mutex.Lock()
	defer campaignCache.mutex.Unlock()
	campaignCache.CampaignsDetails[key] = value
}

func getCampaignDetails(key string) Campaign {
	campaignCache.mutex.RLock()
	defer campaignCache.mutex.RUnlock()
	return campaignCache.CampaignsDetails[key]
}

func (c *CampaignCache) GetCampaignsByCountry(country string) ([]Campaign, error) {
	var result []Campaign
	for _, cid := range getCampaign(country) {
		result = append(result, getCampaignDetails(cid))
	}
	return result, nil
}
