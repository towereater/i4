package db

import (
	"context"
	"fmt"
	"i4-lib/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCollection(ctx context.Context, cfg config.DBConfig, coll string) (*mongo.Collection, error) {
	// Connect to the db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.Host))
	if err != nil {
		return nil, err
	}

	// Retrieve the collection
	return client.Database(cfg.DBName).Collection(coll), nil
}

func getClientCollection(ctx context.Context, cfg config.DBConfig, code string, coll string) (*mongo.Collection, error) {
	// Connect to the db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.Host))
	if err != nil {
		return nil, err
	}

	// Get client db name
	clientDB := getClientDBName(cfg.DBName, code)

	// Retrieve the collection
	return client.Database(getClientDBName(clientDB, code)).Collection(coll), nil
}

func getContextFromConfig(cfg config.DBConfig) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
}

func getClientDBName(db string, code string) string {
	return fmt.Sprintf("%s-%s", db, code)
}
