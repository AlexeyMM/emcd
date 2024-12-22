package app

import (
	"context"
	"fmt"
	"time"
)

func Ticker(job WorkerFn, interval time.Duration) Worker {
	return WorkerFn(func(ctx context.Context) error {
		tick := time.Tick(interval)
		for {
			if err := job.Run(ctx); err != nil {
				return fmt.Errorf("ticker job error: %w", err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tick:
			}
		}
	})
}
