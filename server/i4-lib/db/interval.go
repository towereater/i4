package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SelectInterval(cfg config.DBConfig, client string, dataFilter model.DataInterval, tsFrom string, tsTo string, limit int) ([]model.DataInterval, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Interval)
	if err != nil {
		return []model.DataInterval{}, err
	}

	// Setup find options
	var opts options.FindOptions
	opts.SetLimit(int64(limit))

	// Setup filter
	filter := bson.M{}
	if dataFilter.Machine != "" {
		filter["machine"] = dataFilter.Machine
	}
	if dataFilter.Key != "" {
		filter["key"] = dataFilter.Key
	}
	if dataFilter.Value != "" {
		filter["value"] = dataFilter.Value
	}

	tsFilter := bson.M{}
	if tsFrom != "" {
		tsFilter["gt"] = tsFrom
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}
	if tsTo != "" {
		tsFilter["lt"] = tsTo
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}

	// Define the cursor
	cursor, err := coll.Find(ctx, filter, &opts)
	if err != nil {
		return []model.DataInterval{}, err
	}

	// Search for the documents
	var data []model.DataInterval
	err = cursor.All(ctx, &data)

	return data, err
}

func CountInterval(cfg config.DBConfig, client string, dataFilter model.DataInterval, tsFrom string, tsTo string) (int64, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Interval)
	if err != nil {
		return -1, err
	}

	// Setup filter
	filter := bson.M{}
	if dataFilter.Machine != "" {
		filter["machine"] = dataFilter.Machine
	}
	if dataFilter.Key != "" {
		filter["key"] = dataFilter.Key
	}
	if dataFilter.Value != "" {
		filter["value"] = dataFilter.Value
	}

	tsFilter := bson.M{}
	if tsFrom != "" {
		tsFilter["gt"] = tsFrom
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}
	if tsTo != "" {
		tsFilter["lt"] = tsTo
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}

	// Count the documents
	count, err := coll.CountDocuments(ctx, filter, nil)
	if err != nil {
		return -1, err
	}

	return count, err
}

func SumInterval(cfg config.DBConfig, client string, dataFilter model.DataInterval, tsFrom string, tsTo string) ([]model.DataIntervalSum, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Interval)
	if err != nil {
		return nil, err
	}

	// Setup filter
	filter := bson.M{}
	if dataFilter.Machine != "" {
		filter["machine"] = dataFilter.Machine
	}
	if dataFilter.Key != "" {
		filter["key"] = dataFilter.Key
	}
	if dataFilter.Value != "" {
		filter["value"] = dataFilter.Value
	}

	tsFilter := bson.M{}
	if tsFrom != "" {
		tsFilter["gt"] = tsFrom
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}
	if tsTo != "" {
		tsFilter["lt"] = tsTo
		filter["start"] = tsFilter
		filter["end"] = tsFilter
	}

	// Setup pipeline
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"machine": "$machine",
				"key":     "$key",
				"value":   "$value",
			},
			"count": bson.M{"$sum": 1},
			"sum": bson.M{"$sum": bson.M{"$divide": bson.A{
				bson.M{"$subtract": bson.A{
					bson.M{"$toDate": "$end"},
					bson.M{"$toDate": "$start"},
				}},
				1000,
			}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "machine", Value: 1},
			{Key: "key", Value: 1},
			{Key: "value", Value: 1},
		}}},
	}

	// Run the pipeline
	cursor, err := coll.Aggregate(ctx, pipeline, nil)
	if err != nil {
		return nil, err
	}

	// Search for the documents
	var data []model.DataIntervalSum
	err = cursor.All(ctx, &data)

	return data, err
}

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
