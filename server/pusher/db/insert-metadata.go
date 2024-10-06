package db

import (
	"context"
	"pusher/model"
)

func InsertMetadata(ctx context.Context, metadata model.FileMetadata) error {
	// Retrieve the collection
	coll, err := getCollection(ctx, "i4", "files-content")
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, metadata)

	return err
}
