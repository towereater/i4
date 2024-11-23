package db

import (
	"collector/config"
	"collector/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func SelectMetadata(ctx context.Context, hash uint32) (*model.UploadMetadata, error) {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Metadata)
	if err != nil {
		return nil, err
	}

	// Search for a document
	var metadata model.UploadMetadata
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&metadata)

	return &metadata, err
}
