package service

import (
	"i4-lib/config"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func QueueContent(cfg config.QueueConfig, client string, hash string) error {
	// Create topic writer
	w := &kafka.Writer{
		Addr:  kafka.TCP(cfg.Host),
		Topic: cfg.Topic,
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
