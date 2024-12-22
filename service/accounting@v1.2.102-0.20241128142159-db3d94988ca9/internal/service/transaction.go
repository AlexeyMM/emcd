package service

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"time"
)

type Transaction interface {
	ListTransactions(
		ctx context.Context,
		receiverAccountIds []int,
		types []int,
		from, to time.Time,
		limit int,
		fromTransactionID int64,
	) ([]*model.Transaction, int64, error)
}

type transaction struct {
	transactionRepo repository.Transaction
}

func NewTransaction(
	transactionRepo repository.Transaction,
) Transaction {
	return &transaction{
		transactionRepo: transactionRepo,
	}
}

func (s *transaction) ListTransactions(
	ctx context.Context,
	receiverAccountIds []int,
	types []int,
	from, to time.Time,
	limit int,
	fromTransactionID int64,
) ([]*model.Transaction, int64, error) {
	return s.transactionRepo.ListTransactions(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)
}
