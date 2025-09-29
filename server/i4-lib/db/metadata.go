package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SelectMetadata(cfg config.DBConfig, client string, hash string) (model.UploadMetadata, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Metadata)
	if err != nil {
		return model.UploadMetadata{}, err
	}

	// Search the document
	var metadata model.UploadMetadata
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&metadata)

	return metadata, err
}

func UpsertMetadata(cfg config.DBConfig, client string, metadata model.UploadMetadata) (bool, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Metadata)
	if err != nil {
		return false, err
	}

	// Upsert the document
	result, err := coll.UpdateOne(ctx,
		bson.M{"hash": metadata.Hash},
		bson.M{"$set": metadata},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}
