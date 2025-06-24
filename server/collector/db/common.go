package db

import (
	"collector/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCollection(ctx context.Context, db string, coll string) (*mongo.Collection, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Setup timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.DB.Timeout)*time.Second)
	defer cancel()

	// Connect to the db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.DB.Host))
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

func getClientCollection(ctx context.Context, db string, code string, coll string) (*mongo.Collection, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Setup timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.DB.Timeout)*time.Second)
	defer cancel()

	// Connect to the db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.DB.Host))
	if err != nil {
		return nil, err
	}

	// Check db status
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the collection
	return client.Database(getClientDBName(db, code)).Collection(coll), nil
}

func getClientDBName(db string, code string) string {
	return fmt.Sprintf("%s-%s", db, code)
}
