package db

import (
	"context"
)

func InsertFile(ctx context.Context, user any) error {
	// Retrieve the collection
	coll, err := getCollection(ctx, "i4", "files-content")
	if err != nil {
		return err
	}

	// Insert of a document
	_, err = coll.InsertOne(ctx, user)

	return err
}
