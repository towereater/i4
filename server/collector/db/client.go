package db

import (
	"collector/config"
	"collector/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
