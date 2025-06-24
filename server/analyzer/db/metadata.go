package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

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
