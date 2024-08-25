package database

import (
	"context"

	"github.com/imabg/sync/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseCtx struct {
	context context.Context
	logger config.Application
	uri string
}

func NewDB(ctx context.Context, logger config.Application, uri string) *DatabaseCtx {
	return &DatabaseCtx{
		context: ctx,
		logger: logger,
		uri: uri,
	}
}

// CreateMongoConnection accepts MONGO_URI check the connection & return client pointer
func (dbCtx DatabaseCtx)CreateMongoConnection() (*mongo.Client, error) {
	
	client, err := mongo.Connect(dbCtx.context, options.Client().ApplyURI(dbCtx.uri))
	if err != nil {
		return  nil, err
	}
	dbCtx.logger.InfoLog.Info("Database is connected")
	err = client.Ping(dbCtx.context, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

// DisconnectMongoConnection closes the active connection
func (dbCtx DatabaseCtx)DiscountMongoConnection(client*mongo.Client) error {
	return client.Disconnect(dbCtx.context)
}