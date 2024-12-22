package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Withdraw interface {
	GetBySwapID(ctx context.Context, swapID uuid.UUID) (*model.Withdraw, error)
}

type withdraw struct {
	repo        repository.Withdraw
	explorerRep repository.Explorer
}

func NewWithdraw(repo repository.Withdraw, explorerRep repository.Explorer) *withdraw {
	return &withdraw{
		repo:        repo,
		explorerRep: explorerRep,
	}
}

func (s *withdraw) GetBySwapID(ctx context.Context, swapID uuid.UUID) (*model.Withdraw, error) {
	w, err := s.repo.FindOne(ctx, &model.WithdrawFilter{
		SwapID: &swapID,
	})
	if err != nil {
		return nil, fmt.Errorf("findOne: %w", err)
	}

	var explorerLink string
	if w.HashID != "" {
		explorerLink, err = s.explorerRep.GetTransactionLink(ctx, w.Coin, w.HashID)
		if err != nil {
			log.Error(ctx, "findOneWithHashID: getTransactionLink: %s", err.Error())
		}
	}
	w.ExplorerLink = explorerLink

	return w, nil
}
