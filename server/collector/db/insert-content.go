package db

import (
	"collector/config"
	"collector/model"
	"context"
)

func InsertFile(ctx context.Context, content model.UploadContent) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Content)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, content)

	return err
}
