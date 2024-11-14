package db

import (
	"context"
	"worker/config"
	"worker/model"
)

func InsertInterval(ctx context.Context, content model.Interval) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Interval.Name)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, content)

	return err
}
