package db

import (
	"i4-lib/config"
	"i4-lib/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SelectGauge(cfg config.DBConfig, client string, dataFilter model.DataGauge, tsFrom string, tsTo string, limit int) ([]model.DataGauge, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Gauge)
	if err != nil {
		return []model.DataGauge{}, err
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

	if dataFilter.Timestamp != "" {
		filter["ts"] = dataFilter.Timestamp
	} else {
		tsFilter := bson.M{}
		if tsFrom != "" {
			tsFilter["gt"] = tsFrom
			filter["ts"] = tsFilter
		}
		if tsTo != "" {
			tsFilter["lt"] = tsTo
			filter["ts"] = tsFilter
		}
	}

	// Define the cursor
	cursor, err := coll.Find(ctx, filter, &opts)
	if err != nil {
		return []model.DataGauge{}, err
	}

	// Search for the documents
	var data []model.DataGauge
	err = cursor.All(ctx, &data)

	return data, err
}

func CountGauge(cfg config.DBConfig, client string, dataFilter model.DataGauge, tsFrom string, tsTo string) (int64, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Gauge)
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

	if dataFilter.Timestamp != "" {
		filter["ts"] = dataFilter.Timestamp
	} else {
		tsFilter := bson.M{}
		if tsFrom != "" {
			tsFilter["gt"] = tsFrom
			filter["ts"] = tsFilter
		}
		if tsTo != "" {
			tsFilter["lt"] = tsTo
			filter["ts"] = tsFilter
		}
	}

	// Count the documents
	count, err := coll.CountDocuments(ctx, filter, nil)
	if err != nil {
		return -1, err
	}

	return count, err
}

func SumGauge(cfg config.DBConfig, client string, dataFilter model.DataGauge, tsFrom string, tsTo string) ([]model.DataGaugeSum, error) {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Gauge)
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

	if dataFilter.Timestamp != "" {
		filter["ts"] = dataFilter.Timestamp
	} else {
		tsFilter := bson.M{}
		if tsFrom != "" {
			tsFilter["gt"] = tsFrom
			filter["ts"] = tsFilter
		}
		if tsTo != "" {
			tsFilter["lt"] = tsTo
			filter["ts"] = tsFilter
		}
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
			"sum":   bson.M{"$sum": "$value"},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":     0,
			"machine": "$_id.machine",
			"key":     "$_id.key",
			"value":   "$_id.value",
			"count":   1,
			"sum":     1,
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
	var data []model.DataGaugeSum
	err = cursor.All(ctx, &data)

	return data, err
}

func InsertGauge(cfg config.DBConfig, client string, data model.DataGauge) error {
	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Retrieve the collection
	coll, err := getClientCollection(ctx, cfg, client, cfg.Collections.Gauge)
	if err != nil {
		return err
	}

	// Insert the document
	_, err = coll.InsertOne(ctx, data)

	return err
}
