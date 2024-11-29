package utils

import (
	"analyzer/config"
	"context"
	"encoding/binary"

	"github.com/segmentio/kafka-go"
)

func UnqueueContent(ctx context.Context) (uint32, string, error) {
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
		return 0, "", err
	}

	// Convert the read value
	hash := binary.LittleEndian.Uint32(m.Value[0:4])
	client := string(m.Value[4:8])

	return hash, client, nil
}
