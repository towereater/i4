package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func SelectContent(ctx context.Context, hash uint32) (*model.UploadContent, error) {
	// Extract configuration from context
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Interval)
	if err != nil {
		return nil, err
	}

	// Search for a document
	var content model.UploadContent
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&content)

	return &content, err
}
