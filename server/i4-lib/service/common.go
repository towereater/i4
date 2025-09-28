package service

import (
	"context"
	"i4-lib/config"
	"time"
)

func getContextFromConfig(cfg config.QueueConfig) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
}
