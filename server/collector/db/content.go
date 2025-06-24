package db

import (
	"collector/config"
	"collector/model"
	"context"
)

func InsertContent(ctx context.Context, client string, content model.UploadContent) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg.DB.DBName, client, cfg.DB.Collections.Content)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, content)

	return err
}
