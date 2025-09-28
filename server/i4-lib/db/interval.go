package db

import (
	"i4-lib/config"
	"i4-lib/model"
)

func InsertInterval(cfg config.DBConfig, client string, data model.DataInterval) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Interval)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, data)

	return err
}
