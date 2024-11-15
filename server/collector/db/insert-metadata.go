package db

import (
	"collector/config"
	"collector/model"
	"context"
)

func InsertMetadata(ctx context.Context, metadata model.UploadMetadata) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Metadata)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}
