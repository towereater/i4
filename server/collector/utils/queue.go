package utils

import (
	"collector/config"
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func QueueContent(ctx context.Context, client string, hash string) error {
	// Extract config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Create topic writer
	w := &kafka.Writer{
		Addr:  kafka.TCP(cfg.Queue.Host),
		Topic: cfg.Queue.Topic,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Queue.Timeout)*time.Second)
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
