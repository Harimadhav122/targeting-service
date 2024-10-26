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
	mongo_conn_uri = "mongodb://localhost:27017/"
)

type IMongo interface {
	GetDb(db_name string) IMongoDb
}

type IMongoDb interface {
	GetCollection(coll_name string) IMongoCollection
}

type IMongoCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error)
	Aggregate(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

type Mongo struct {
}

type MongoDb struct {
	Db *mongo.Database
}

type MongoCollection struct {
	Collection *mongo.Collection
}

var logger log.Logger
var MongoDB IMongo

func init() {

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "mongodb")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

	conn_uri := os.Getenv("MONGODB_CONN_URI")
	if conn_uri != "" {
		mongo_conn_uri = conn_uri
	}
	MongoDB = Mongo{}
}

func (m Mongo) GetDb(db_name string) IMongoDb {
	level.Info(logger).Log("msg", "Attempting to connect to mongodb..", "connection uri", mongo_conn_uri)
	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_conn_uri))
	if err != nil {
		level.Error(logger).Log("method", "GetDb", "err", err)
		panic(err)
	}
	level.Info(logger).Log("msg", "Connected to mongodb succesfully")
	return &MongoDb{Db: client.Database(db_name)}
}

func (m *MongoDb) GetCollection(coll_name string) IMongoCollection {
	return &MongoCollection{Collection: m.Db.Collection(coll_name)}
}

func (m *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	doc := m.Collection.FindOne(ctx, filter, opts...)
	return doc, nil
}

func (m *MongoCollection) Aggregate(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := m.Collection.Aggregate(ctx, filter, opts...)
	if err != nil {
		level.Error(logger).Log("msg", "mongodb aggregate failed", "err", err)
		return nil, err
	}
	return cursor, nil
}
