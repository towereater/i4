package db

import (
	"collector/config"
	"collector/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertMetadata(ctx context.Context, client string, metadata model.UploadMetadata) (bool, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg.DB.DBName, client, cfg.DB.Collections.Metadata)
	if err != nil {
		return false, err
	}

	// Upsert the document
	result, err := coll.UpdateOne(ctx,
		bson.M{"hash": metadata.Hash},
		bson.M{"$set": metadata},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}

func SelectMetadata(ctx context.Context, client string, hash string) (*model.UploadMetadata, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg.DB.DBName, client, cfg.DB.Collections.Metadata)
	if err != nil {
		return nil, err
	}

	// Search the document
	var metadata model.UploadMetadata
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&metadata)

	return &metadata, err
}
