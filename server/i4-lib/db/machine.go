package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertMachine(cfg config.DBConfig, client string, machine model.Machine) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return err
	}

	// Update the document
	_, err = coll.UpdateOne(ctx,
		bson.M{"code": client, "machines.code": bson.M{"$ne": machine.Code}},
		bson.M{"$push": bson.M{"machines": machine}},
	)

	return err
}

func RemoveMachine(cfg config.DBConfig, client string, machine string) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getCollection(ctx, cfg, cfg.Collections.Clients)
	if err != nil {
		return err
	}

	// Update the document
	_, err = coll.UpdateOne(ctx,
		bson.M{"code": client},
		bson.M{"$pull": bson.M{"machines": bson.M{"code": machine}}},
	)

	return err
}
