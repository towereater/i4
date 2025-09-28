package db

import (
	"i4-lib/config"
	"i4-lib/model"
)

func InsertContent(cfg config.DBConfig, client string, content model.UploadContent) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Content)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, content)

	return err
}
