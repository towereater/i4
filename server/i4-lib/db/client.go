package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SelectClientByCode(cfg config.DBConfig, code string) (model.Client, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return model.Client{}, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{"code": code}).Decode(&client)

	return client, err
}

func SelectClientByApiKey(cfg config.DBConfig, apiKey string) (model.Client, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return model.Client{}, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{"apiKey": apiKey}).Decode(&client)

	return client, err
}

func SelectAnyClient(cfg config.DBConfig) (model.Client, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return model.Client{}, err
	}

	// Search the document
	var client model.Client
	err = coll.FindOne(ctx, bson.M{}).Decode(&client)

	return client, err
}

func InsertClient(cfg config.DBConfig, client model.Client) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, client)

	return err
}

func SetupClientCollections(cfg config.DBConfig, code string) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the metadata collection
	coll, err := getClientCollection(ctx, cfg, code, cfg.Collections.Metadata)
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

	// Retrieve the metadata collection
	coll, err = getClientCollection(ctx, cfg, code, cfg.Collections.Content)
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

	// Retrieve the metadata collection
	coll, err = getClientCollection(ctx, cfg, code, cfg.Collections.Gauge)
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

	// Retrieve the metadata collection
	coll, err = getClientCollection(ctx, cfg, code, cfg.Collections.Interval)
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

	return nil
}
