package service

import (
	"context"

	"github.com/google/uuid"
)

type DepositAddressPool interface {
	// GetOrCreate probably, should run in a repeatable-read tx, so address doesn't get selected twice
	GetOrCreate(
		ctx context.Context,
		merchantID uuid.UUID,
		networkID string,
		idempotencyKey uuid.UUID,
	) (string, error)
}
