package db

import (
	"collector/config"
	"collector/model"
	"context"
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
