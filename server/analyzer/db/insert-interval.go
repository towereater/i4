package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"
	"fmt"
)

func InsertInterval(ctx context.Context, client string, data model.DataInterval) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx,
		fmt.Sprintf("%s-%s", cfg.DB.DBName, client),
		cfg.DB.Collections.Interval)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, data)

	return err
}
