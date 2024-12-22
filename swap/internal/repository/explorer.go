package repository

import "context"

type Explorer interface {
	GetTransactionLink(ctx context.Context, coin, hashID string) (string, error)
}
