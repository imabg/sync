package database

import (
	"context"

	"github.com/imabg/sync/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseCtx struct {
	context context.Context
	config  config.Application
	log config.Logger
}

func NewDB(ctx context.Context, config config.Application) *DatabaseCtx {
	return &DatabaseCtx{
		context: ctx,
		config:  config,
		log: config.Log,
	}
}

// CreateMongoConnection accepts MONGO_URI check the connection & return client pointer
func (dbCtx DatabaseCtx) CreateMongoConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(dbCtx.context, options.Client().ApplyURI(dbCtx.config.Env.MongoURI))
	if err != nil {
		return nil, err
	}
	err = client.Ping(dbCtx.context, readpref.Primary())
	if err != nil {
		return nil, err
	}
	dbCtx.log.InfoLog.Info("Database is connected")
	return client, err
}

// DisconnectMongoConnection closes the active connection
func (dbCtx DatabaseCtx) DiscountMongoConnection(client *mongo.Client) error {
	return client.Disconnect(dbCtx.context)
}

func (dbCtx DatabaseCtx) GetMongoDatabase(client *mongo.Client) *mongo.Database {
	return client.Database(dbCtx.config.Env.DBName)
}