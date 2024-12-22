package worker

import (
	"context"
	"time"

	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/emcd/sdk/log"
)

//go:generate moq -out coin_fetcher_repository_moq_test.go . Repository

type Repository interface {
	FetchCoins(ctx context.Context) error
}

type CoinFetcher struct {
	coinService service.CoinService

	fetchFrequency time.Duration
}

func NewCoinFetcher(coinService service.CoinService, fetchFrequency time.Duration) *CoinFetcher {
	return &CoinFetcher{coinService: coinService, fetchFrequency: fetchFrequency}
}

func (c *CoinFetcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.fetchFrequency)
	defer ticker.Stop()

	// Fetch immediately after starting, then fetch at the defined frequency.
	// This approach addresses the issue of having an empty cache in the initial few minutes
	// if the fetch frequency is set to a few minutes.
	if err := c.coinService.UpdateCoins(ctx); err != nil {
		log.SError(ctx, "fetch coins error", map[string]any{"error": err})
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := c.coinService.UpdateCoins(ctx); err != nil {
				log.SError(ctx, "fetch coins error", map[string]any{"error": err})
			}
		}
	}
}
