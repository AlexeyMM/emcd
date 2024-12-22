package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/repository"
)

type SwapCoin interface {
	GetAll(ctx context.Context) ([]*model.Coin, error)
}

type swapCoin struct {
	repo repository.Coin
}

func NewSwapCoin(repo repository.Coin) *swapCoin {
	return &swapCoin{
		repo: repo,
	}
}

func (s *swapCoin) GetAll(ctx context.Context) ([]*model.Coin, error) {
	coins, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAll: %w", err)
	}

	return coins, nil
}
