package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"analyzer/config"
	"analyzer/db"
	"analyzer/model"
	"analyzer/utils"
)

func main() {
	// Get run args
	if len(os.Args) < 2 {
		fmt.Printf("No config file set\n")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Printf("Loading configuration from %s\n", configPath)
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Error while reading config file: %s\n", err.Error())
		os.Exit(2)
	}
	ctx := context.WithValue(context.Background(), config.ContextConfig, cfg)

	// Main loop
	for {
		// Poll the queue for data
		client, hash, err := utils.UnqueueContent(ctx)
		fmt.Printf("Received client %s and hash %s\n", client, hash)
		if err != nil {
			fmt.Printf("Error while reading queue: %s\n", err.Error())
			continue
		}

		// Get content from db
		content, err := db.SelectContent(ctx, client, hash)
		if err != nil {
			fmt.Printf("Error while reading content from db: %s\n", err.Error())
			continue
		}

		fmt.Printf("File content is:\n%s\n", string(content.Content))

		// Convert and split the content
		data := strings.Split(string(content.Content), "\n")

		// Save data to db
		for _, r := range data {
			if r == "" {
				continue
			}

			// Convert the json to struct
			var dataContent model.DataContent
			err = json.Unmarshal([]byte(r), &dataContent)
			if err != nil {
				fmt.Printf("Error while converting data content to string: %s\n", err.Error())
				continue
			}

			jsonString, err := json.Marshal(dataContent.Content)
			if err != nil {
				fmt.Printf("Error while converting content to json: %s\n", err.Error())
				continue
			}

			switch dataContent.Type {
			// Gauge data
			case "GAU":
				var gauge model.DataGauge
				err = json.Unmarshal(jsonString, &gauge)
				if err != nil {
					fmt.Printf("Error while converting data gauge to json: %s\n", err.Error())
					continue
				}

				err = db.InsertGauge(ctx, client, gauge)
			// Interval data
			case "INT":
				var interval model.DataInterval
				err = json.Unmarshal(jsonString, &interval)
				if err != nil {
					fmt.Printf("Error while converting data gauge to json: %s\n", err.Error())
					continue
				}

				err = db.InsertInterval(ctx, client, interval)
			default:
				err = fmt.Errorf("undefined data type")
			}
			if err != nil {
				fmt.Printf("Error while converting content to string: %s\n", err.Error())
			}
		}
	}
}
