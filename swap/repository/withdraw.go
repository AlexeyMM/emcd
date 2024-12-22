package repository

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/protocol/swapWithdraw"
	"github.com/google/uuid"
)

type Withdraw interface {
	GetTransactionLink(ctx context.Context, swapID uuid.UUID) (string, error)
}

type withdraw struct {
	handler swapWithdraw.SwapWithdrawServiceClient
}

func NewWithdraw(handler swapWithdraw.SwapWithdrawServiceClient) *withdraw {
	return &withdraw{
		handler: handler,
	}
}

func (w *withdraw) GetTransactionLink(ctx context.Context, swapID uuid.UUID) (string, error) {
	resp, err := w.handler.GetTransactionLink(ctx, &swapWithdraw.GetTransactionLinkRequest{
		SwapId: swapID.String(),
	})
	if err != nil {
		return "", fmt.Errorf("getTransactionLink: %w", err)
	}
	return resp.Link, nil
}
