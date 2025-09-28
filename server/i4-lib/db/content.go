package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
)

func SelectContent(cfg config.DBConfig, client string, hash string) (model.UploadContent, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Content)
	if err != nil {
		return model.UploadContent{}, err
	}

	// Search the document
	var content model.UploadContent
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&content)

	return content, err
}

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
