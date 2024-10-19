package storage

import (
	"context"
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	supported_countries = [4]string{"us", "india", "canada", "germany"}
	supported_os        = [2]string{"android", "ios"}
	supported_apps      = [5]string{"com.abc.xyz", "com.gametion.ludokinggame", "com.apple.in", "com.google.in", "com.samsung.in"}
	mongo_conn_uri      = "mongodb://localhost:27017/campaigns"
	db_name             = "campaigns"
	collection_name     = "campaigns_details"
	ctx                 = context.Background()
	logger              log.Logger
	wg                  sync.WaitGroup
)

type ICampaignCache interface {
	GetCampaignsByCountry(country string) ([]Campaign, error)
}

type CampaignCache struct {
	mutex            sync.RWMutex
	Campaigns        map[string][]string // country -> campaigns
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

type CampaignIdResponse struct {
	Id string `bson:"_id"`
}

type CampaignResponse struct {
	Id             string       `bson:"_id"`
	Name           string       `bson:"name"`
	Image          string       `bson:"image"`
	Cta            string       `bson:"cta"`
	IsActive       bool         `bson:"isActive"`
	NoRestrictions bool         `bson:"noRestrictions"`
	Rules          RuleResponse `bson:"rules"`
}

type RuleResponse struct {
	IncludeOs  []string `bson:"includeOs,omitempty"`
	ExcludeOs  []string `bson:"excludeOs,omitempty"`
	IncludeApp []string `bson:"includeApp,omitempty"`
	ExcludeApp []string `bson:"excludeApp,omitempty"`
}

var campaignCache = &CampaignCache{Campaigns: make(map[string][]string), CampaignsDetails: make(map[string]Campaign)}

func init() {

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "cache")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

	level.Info(logger).Log("method", "getAndFillCache", "msg", "cache fill started")
	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_conn_uri))
	if err != nil {
		level.Error(logger).Log("method", "getAndFillCache", "err", err)
		panic(err)
	}
	db := client.Database(db_name)

	// increment goroutine counter
	wg.Add(1)
	// fill cache with all campaigns details
	go campaignCache.getCampaignsAndFillCache(db)

	// fill cache with campaigns per country
	for _, country := range supported_countries {
		wg.Add(1)
		go campaignCache.getCampaignsByCountryAndFillCache(db, country)
	}
	// wait until all goroutines have finished
	wg.Wait()
	level.Info(logger).Log("method", "getAndFillCache", "msg", "cache fill completed")
}

func (c *CampaignCache) getCampaignsAndFillCache(db *mongo.Database) {
	// goroutine execution completed
	defer wg.Done()
	var data []CampaignResponse
	cursor, err := db.Collection(collection_name).Find(ctx, bson.D{})
	if err != nil {
		level.Error(logger).Log("method", "getAndFillCache", "err", err)
	}

	if err = cursor.All(context.TODO(), &data); err != nil {
		level.Error(logger).Log("method", "getAndFillCache", err, "error decoding cursor")
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
		c.setCampaignDetails(d.Id, campaign)
	}
}

func (c *CampaignCache) getCampaignsByCountryAndFillCache(db *mongo.Database, country string) {
	// goroutine execution completed
	defer wg.Done()
	cursor, err := db.Collection(country).Find(ctx, bson.D{})
	if err != nil {
		level.Error(logger).Log("method", "getAndFillCache", "err", err)
	}
	var result []string
	var data []CampaignIdResponse

	if err = cursor.All(context.TODO(), &data); err != nil {
		level.Error(logger).Log("method", "getAndFillCache", err, "error decoding cursor")
		panic(err)
	}
	for _, d := range data {
		result = append(result, d.Id)
	}
	c.setCampaign(country, result)
}

func NewCache() ICampaignCache {
	return campaignCache
}

func (c *CampaignCache) setCampaign(key string, value []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Campaigns[key] = value
}

func (c *CampaignCache) getCampaign(key string) []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Campaigns[key]
}

func (c *CampaignCache) setCampaignDetails(key string, value Campaign) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.CampaignsDetails[key] = value
}

func (c *CampaignCache) getCampaignDetails(key string) Campaign {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.CampaignsDetails[key]
}

func (c *CampaignCache) GetCampaignsByCountry(country string) ([]Campaign, error) {
	var result []Campaign
	for _, cid := range c.getCampaign(country) {
		result = append(result, c.getCampaignDetails(cid))
	}
	return result, nil
}
