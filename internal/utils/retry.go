package utils

import (
	"context"
	"fmt"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	RetryDelay time.Duration
}

func WithRetry(ctx context.Context, fn func() error, config RetryConfig) error {
	var lastErr error
	for i := 0; i <= config.MaxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if i == config.MaxRetries {
				break
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(config.RetryDelay):
				continue
			}
		}
		return nil
	}
	return fmt.Errorf("failed after %d retries: %v", config.MaxRetries, lastErr)
}
