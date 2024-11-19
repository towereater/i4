package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"
)

func InsertInterval(ctx context.Context, data model.DataInterval) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Interval)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, data)

	return err
}
