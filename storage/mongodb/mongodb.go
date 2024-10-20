package mongodb

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongo_conn_uri = "mongodb://localhost:27017/campaigns"
	db_name        = "campaigns"
)

type IMongo interface {
	Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	Aggregate(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

type Mongo struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type CampaignResponse struct {
	Id             string       `bson:"_id"`
	Name           string       `bson:"name"`
	Image          string       `bson:"image"`
	Cta            string       `bson:"cta"`
	IsActive       bool         `bson:"isActive"`
	NoRestrictions bool         `bson:"noRestrictions"`
	Rules          RuleResponse `bson:"rules,omitempty"`
}

type RuleResponse struct {
	IncludeOs  []string `bson:"includeOs,omitempty"`
	ExcludeOs  []string `bson:"excludeOs,omitempty"`
	IncludeApp []string `bson:"includeApp,omitempty"`
	ExcludeApp []string `bson:"excludeApp,omitempty"`
}

type CampaignIdResponse struct {
	Id string `bson:"_id"`
}

var logger log.Logger
var mongoInstance = &Mongo{}

func init() {

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "mongodb")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

	level.Info(logger).Log("msg", "Attempting to connect to mongodb..")
	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_conn_uri))
	if err != nil {
		level.Error(logger).Log("method", "init", "err", err)
		panic(err)
	}
	level.Info(logger).Log("msg", "Connected to mongodb succesfully")
	db := client.Database(db_name)
	mongoInstance = &Mongo{
		Client: client,
		Db:     db,
	}
}

func NewMongo() IMongo {
	return mongoInstance
}

func (m *Mongo) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := m.Db.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		level.Error(logger).Log("msg", "mongodb find failed", "err", err)
		return nil, err
	}
	return cursor, nil
}

func (m *Mongo) Aggregate(ctx context.Context, collection string, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := m.Db.Collection(collection).Aggregate(ctx, filter, opts...)
	if err != nil {
		level.Error(logger).Log("msg", "mongodb aggregate failed", "err", err)
		return nil, err
	}
	return cursor, nil
}
