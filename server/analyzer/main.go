package main

import (
	"encoding/json"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Get run args
	if len(os.Args) < 2 {
		service.Log("No config file set\n")
		os.Exit(1)
	}
	configPath := os.Args[1]
	service.Log("Loading configuration from %s\n", configPath)

	// Setup machine config
	var cfg config.BaseConfig
	err := config.LoadConfig(configPath, &cfg)
	if err != nil {
		service.Log("Error while reading config file: %s\n", err.Error())
		os.Exit(2)
	}
	service.Log("Configuration loaded: %+v\n", cfg)

	// Main loop
	for {
		// Poll the queue for data
		client, hash, err := service.UnqueueContent(cfg.Queue)
		service.Log("Received client %s and hash %s\n", client, hash)
		if err != nil {
			service.Log("Error while reading queue: %s\n", err.Error())
			continue
		}

		// Get content from db
		content, err := db.SelectContent(cfg.DB, client, hash)
		if err != nil {
			service.Log("Error while reading content from db: %s\n", err.Error())
			continue
		}

		// Convert and split the content
		data := strings.Split(string(content.Content), "\n")

		// Save data to db
		for _, r := range data {
			if r == "" {
				service.Log("Row skipped because empty data")
				continue
			}

			// Convert the json to struct
			var dataContent model.DataContent
			err = json.Unmarshal([]byte(r), &dataContent)
			if err != nil {
				service.Log("Error while converting common data content %s: %s\n", r, err.Error())
				continue
			}

			jsonString, err := json.Marshal(dataContent.Content)
			if err != nil {
				service.Log("Error while converting specific data content %+v: %s\n", dataContent.Content, err.Error())
				continue
			}

			switch dataContent.Type {
			// Gauge data
			case "GAU":
				var gauge model.DataGauge
				err = json.Unmarshal(jsonString, &gauge)
				if err != nil {
					service.Log("Error while converting data gauge %s: %s\n", jsonString, err.Error())
					continue
				}

				err = db.InsertGauge(cfg.DB, client, gauge)
				if mongo.IsDuplicateKeyError(err) {
					service.Log("Error while inserting data gauge because of duplicate key: %+v", gauge)
					continue
				} else if err != nil {
					service.Log("Error while saving data gauge %+v: %s\n", gauge, err.Error())
					continue
				}
			// Interval data
			case "INT":
				var interval model.DataInterval
				err = json.Unmarshal(jsonString, &interval)
				if err != nil {
					service.Log("Error while converting data interval: %s\n", err.Error())
					continue
				}

				err = db.InsertInterval(cfg.DB, client, interval)
				if mongo.IsDuplicateKeyError(err) {
					service.Log("Error while inserting data interval because of duplicate key: %+v", interval)
					continue
				} else if err != nil {
					service.Log("Error while saving data interval %+v: %s\n", interval, err.Error())
					continue
				}
			default:
				service.Log("Undefined data type %s", dataContent.Type)
				continue
			}
		}
	}
}
