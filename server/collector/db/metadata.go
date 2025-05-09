package db

import (
	"collector/config"
	"collector/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertMetadata(ctx context.Context, metadata model.UploadMetadata) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Metadata)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}

func SelectMetadata(ctx context.Context, hash uint32) (*model.UploadMetadata, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Metadata)
	if err != nil {
		return nil, err
	}

	// Search the document
	var metadata model.UploadMetadata
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&metadata)

	return &metadata, err
}
