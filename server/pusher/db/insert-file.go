package db

import (
	"context"
	"pusher/config"
)

func InsertFile(ctx context.Context, user any) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Content.Name)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, user)

	return err
}
