package utils

import (
	"collector/config"
	"context"
	"encoding/binary"
	"time"

	"github.com/segmentio/kafka-go"
)

func QueueContent(ctx context.Context, hash uint32, client string) error {
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
	datetime := time.Now().Format(time.DateTime)
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, hash)
	value := append(h, client...)

	// Write data on queue
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(datetime),
			Value: value,
		},
	)

	return err
}
