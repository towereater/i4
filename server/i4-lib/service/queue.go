package service

import (
	"context"
	"i4-lib/config"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func QueueContent(cfg config.QueueConfig, client string, hash string) error {
	// Create topic writer
	w := &kafka.Writer{
		Addr:  kafka.TCP(cfg.Host),
		Topic: cfg.Topics.Uploads.Topic,
	}

	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

	// Prepare data for queue
	key := []byte(time.Now().Format(time.DateTime))
	value := []byte(strings.Join([]string{client, hash}, ""))

	// Write data on queue
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)

	return err
}

func UnqueueContent(cfg config.QueueConfig) (string, string, error) {
	// Create topic reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.Topics.Uploads.Group,
		Topic:   cfg.Topics.Uploads.Topic,
	})
	defer r.Close()

	// Setup timeout
	ctx, cancel := getContextFromConfig(cfg)
	defer cancel()

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

func getContextFromConfig(cfg config.QueueConfig) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
}
