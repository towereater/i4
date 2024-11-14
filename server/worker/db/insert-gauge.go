package db

import (
	"context"
	"worker/config"
	"worker/model"
)

func InsertGauge(ctx context.Context, metadata model.Gauge) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Gauge.Name)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}
