package db

import (
	"analyzer/config"
	"analyzer/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func SelectContent(ctx context.Context, client string, hash string) (*model.UploadContent, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg.DB.DBName, client, cfg.DB.Collections.Content)
	if err != nil {
		return nil, err
	}

	// Search the document
	var content model.UploadContent
	err = coll.FindOne(ctx, bson.M{"hash": hash}).Decode(&content)

	return &content, err
}
