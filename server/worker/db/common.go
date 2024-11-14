package db

import (
	"context"
	"time"
	"worker/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCollection(ctx context.Context, db string, coll string) (*mongo.Collection, error) {
	// Extracting config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Timeout setup
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.DB.Timeout)*time.Second)
	defer cancel()

	// Connection to the db
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://"+cfg.DB.Host+":"+cfg.DB.Port))
	if err != nil {
		return nil, err
	}

	// Check db status
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the collection
	return client.Database(db).Collection(coll), nil
}
