package db

import (
	"context"
	"pusher/config"
	"pusher/model"
)

func InsertMetadata(ctx context.Context, metadata model.FileMetadata) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Metadata.Name)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}
