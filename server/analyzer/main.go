package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
		println("Error while reading config file:", err)
		os.Exit(2)
	}
	ctx := context.WithValue(context.Background(), config.ContextConfig, cfg)

	// Main loop
	for {
		// Poll the queue for data
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Queue.Timeout)*time.Second)
		hash, client, err := unqueueContent(ctx)
		cancel()
		if err != nil {
			println("Error while reading queued content:", err)
			continue
		}

		// Get content from database
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.DB.Timeout)*time.Second)
		uplContent, err := db.SelectContent(ctx, hash)
		cancel()
		if err != nil {
			println("Error while reading from database:", err)
			continue
		}

		// Convert and split the content
		data := strings.Split(string(uplContent.Content), "\n")
		println(data)

		// Save data to database
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.DB.Timeout)*time.Second)
		for _, r := range data {
			// Parsing of the content
			var content model.DataContent
			err = json.Unmarshal([]byte(r), &content)
			if err != nil {
				println("Error while converting content:", err)
				continue
			}

			switch content.Type {
			// Content is a interval
			case "INT":
				err = db.InsertInterval(ctx, client, content.Content.(model.DataInterval))
			// Content is a gauge
			case "GAU":
				err = db.InsertGauge(ctx, client, content.Content.(model.DataGauge))
			}
			if err != nil {
				println("Error while converting data:", err)
			}
		}
		cancel()
	}
}

func unqueueContent(ctx context.Context) (uint32, string, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Create topic reader with timeout
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Queue.Brokers,
		GroupID:  cfg.Queue.Uploads.Group,
		Topic:    cfg.Queue.Uploads.Topic,
		MaxBytes: 10e6,
	})
	defer r.Close()

	// Read the message from queue
	m, err := r.ReadMessage(ctx)
	if err != nil {
		println("Error while reading queued content:", err)
		os.Exit(3)
	}

	// Convert the read value
	hash := binary.LittleEndian.Uint32(m.Value[0:4])
	client := string(m.Value[4:4])

	return hash, client, r.Close()
}
