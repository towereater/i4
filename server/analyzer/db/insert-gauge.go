package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"
)

func InsertGauge(ctx context.Context, metadata model.DataGauge) error {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Gauge)
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}