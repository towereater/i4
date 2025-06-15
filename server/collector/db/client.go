package db

import (
	"collector/config"
	"collector/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertClient(ctx context.Context, client model.Client) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, client)

	return err
}

func SelectClientByCode(ctx context.Context, code string) (*model.Client, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return nil, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{"code": code}).Decode(&client)

	return &client, err
}

func SelectClientByApiKey(ctx context.Context, apiKey string) (*model.Client, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return nil, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{"apiKey": apiKey}).Decode(&client)

	return &client, err
}

func SelectAnyClient(ctx context.Context) (*model.Client, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return nil, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{}).Decode(&client)

	return &client, err
}

func SetupClientCollections(ctx context.Context, code string) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Construct the database name
	clientDb := fmt.Sprintf("%s-%s", cfg.DB.DBName, code)

	// Retrieve the metadata collection
	coll, err := getCollection(ctx, clientDb, cfg.DB.Collections.Metadata)
	if err != nil {
		return err
	}

	// Create the indexes
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "hash", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	indexModel = mongo.IndexModel{
		Keys: bson.D{
			{Key: "ts", Value: 1},
		},
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Retrieve the content collection
	coll, err = getCollection(ctx, clientDb, cfg.DB.Collections.Content)
	if err != nil {
		return err
	}

	// Create the indexes
	indexModel = mongo.IndexModel{
		Keys: bson.D{
			{Key: "hash", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Retrieve the data gauge collection
	coll, err = getCollection(ctx, clientDb, cfg.DB.Collections.Gauge)
	if err != nil {
		return err
	}

	// Create the indexes
	indexModel = mongo.IndexModel{
		Keys: bson.D{
			{Key: "machine", Value: 1},
			{Key: "key", Value: 1},
			{Key: "ts", Value: 1},
		},
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Retrieve the data interval collection
	coll, err = getCollection(ctx, clientDb, cfg.DB.Collections.Interval)
	if err != nil {
		return err
	}

	// Create the indexes
	indexModel = mongo.IndexModel{
		Keys: bson.D{
			{Key: "machine", Value: 1},
			{Key: "key", Value: 1},
			{Key: "start", Value: 1},
			{Key: "end", Value: 1},
		},
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return err
}
