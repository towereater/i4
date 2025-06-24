package utils

import (
	"analyzer/config"
	"context"

	"github.com/segmentio/kafka-go"
)

func UnqueueContent(ctx context.Context) (string, string, error) {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Create topic reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Queue.Brokers,
		GroupID: cfg.Queue.Uploads.Group,
		Topic:   cfg.Queue.Uploads.Topic,
	})
	defer r.Close()

	// Read the message from queue
	m, err := r.ReadMessage(ctx)
	if err != nil {
		return "", "", err
	}

	// Convert the read value to single components
	client := string(m.Value[0:5])
	hash := string(m.Value[5:69])

	return client, hash, nil
}
