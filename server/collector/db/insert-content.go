package db

import (
	"collector/config"
	"collector/model"
	"context"
)

func InsertContent(ctx context.Context, content model.UploadContent) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Content)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, content)

	return err
}
