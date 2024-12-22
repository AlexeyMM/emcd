package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Transfer interface {
	GetLastInternalTransfer(ctx context.Context, accountID int64) (*model.InternalTransfer, error)
}

type transfer struct {
	repo repository.Transfer
}

func NewTransfer(repo repository.Transfer) *transfer {
	return &transfer{
		repo: repo,
	}
}

func (s *transfer) GetLastInternalTransfer(ctx context.Context, accountID int64) (*model.InternalTransfer, error) {
	isLast := true
	tr, err := s.repo.FindOne(ctx, &model.InternalTransferFilter{
		FromAccountID: &accountID,
		IsLast:        &isLast,
	})
	if err != nil {
		return nil, fmt.Errorf("getLastInternalTransfer: %w", err)
	}
	return tr, nil
}
