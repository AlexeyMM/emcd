package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/repository"
	"github.com/google/uuid"
)

type SwapWithdraw interface {
	GetTransactionLink(ctx context.Context, swapID uuid.UUID) (string, error)
}

type swapWithdraw struct {
	repo repository.Withdraw
}

func NewSwapWithdraw(repo repository.Withdraw) *swapWithdraw {
	return &swapWithdraw{
		repo: repo,
	}
}

func (s *swapWithdraw) GetTransactionLink(ctx context.Context, swapID uuid.UUID) (string, error) {
	link, err := s.repo.GetTransactionLink(ctx, swapID)
	if err != nil {
		return "", fmt.Errorf("getTransactionLink: %w", err)
	}
	return link, nil
}
