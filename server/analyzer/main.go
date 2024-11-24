package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"analyzer/config"
	"analyzer/db"
	"analyzer/model"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Get run args
	if len(os.Args) < 2 {
		println("No config file set")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Println("Loading configuration")
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		println("Error while reading config file:", err.Error())
		os.Exit(2)
	}
	ctx := context.WithValue(context.Background(), config.ContextConfig, cfg)

	// Main loop
	for {
		// Poll the queue for data
		hash, client, err := unqueueContent(ctx)
		if err != nil {
			println("Error while reading queued content:", err.Error())
			continue
		}

		// Get metadata from database
		metadata, err := db.SelectMetadata(ctx, hash)
		if err != nil {
			println("Error while reading from metadata db:", err.Error())
			return
		}

		// Get content from database
		content, err := db.SelectContent(ctx, hash)
		if err != nil {
			println("Error while reading from content db:", err.Error())
			continue
		}

		// Convert and split the content
		data := strings.Split(string(content.Content), "\n")
		fmt.Printf("data is: %+v\n", data)

		// Save data to database
		for _, r := range data {
			if r == "" {
				continue
			}
			fmt.Printf("r is: %+v\n", r)

			// Parsing of the content
			var dataContent model.DataContent
			err = json.Unmarshal([]byte(r), &dataContent)
			if err != nil {
				println("Error while converting data content to string:", err.Error())
				continue
			}

			fmt.Printf("dataContent is: %+v\n", dataContent)
			fmt.Printf("dataContent.Content is: %+v\n", dataContent.Content)
			jsonString, err := json.Marshal(dataContent.Content)
			if err != nil {
				println("Error while converting content to json:", err.Error())
				continue
			}

			switch dataContent.Type {
			// Content is a gauge
			case "GAU":
				var gauge model.DataGauge
				json.Unmarshal(jsonString, &gauge)
				gauge.Machine = metadata.Machine
				err = db.InsertGauge(ctx, client, gauge)
			// Content is a interval
			case "INT":
				var interval model.DataInterval
				json.Unmarshal(jsonString, &interval)
				interval.Machine = metadata.Machine
				err = db.InsertInterval(ctx, client, interval)
			}
			if err != nil {
				println("Error while converting content to string:", err.Error())
			}
		}
	}
}

func unqueueContent(ctx context.Context) (uint32, string, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Create topic reader with timeout
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Queue.Brokers,
		GroupID: cfg.Queue.Uploads.Group,
		Topic:   cfg.Queue.Uploads.Topic,
	})
	defer r.Close()

	// Read the message from queue
	m, err := r.ReadMessage(ctx)
	if err != nil {
		println("Error while reading queued content:", err.Error())
		os.Exit(3)
	}

	// Convert the read value
	hash := binary.LittleEndian.Uint32(m.Value[0:4])
	client := string(m.Value[4:8])

	return hash, client, r.Close()
}
