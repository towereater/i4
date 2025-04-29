package db

import (
	"collector/config"
	"collector/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertMachine(ctx context.Context, client string, machine model.Machine) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.UpdateOne(ctx,
		bson.M{"code": client, "machines.code": bson.M{"$ne": machine.Code}},
		bson.M{"$push": bson.M{"machines": machine}},
	)

	return err
}

func RemoveMachine(ctx context.Context, client string, machine string) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg.DB.DBName, cfg.DB.Collections.Clients)
	if err != nil {
		return err
	}

	// Remove the document
	_, err = coll.UpdateOne(ctx,
		bson.M{"code": client},
		bson.M{"$pull": bson.M{"machines": bson.M{"code": machine}}},
	)

	return err
}
