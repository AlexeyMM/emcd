package service

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Payouts interface {
	FindPayoutsForBlock(ctx context.Context, coinID int64, minPay float32, timestamp time.Time) ([]*model.PayoutForBlock, error)
	GetCurrentPayoutsBlock(ctx context.Context, coinID int64, username string, isService bool) ([]*model.PayoutBlockTransaction, error)
	GetFreePayouts(ctx context.Context, coinID int64) ([]*model.FreePayout, error)
	GetCurrentPayoutsList(ctx context.Context, coinID, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error)
	GetCurrentReferralsPayoutsList(ctx context.Context, coinID, paymentTransactionType, referralId int64) ([]*model.PayoutWithCalculation, error)
	CheckFreePayoutTransaction(ctx context.Context, accountId, transactionId int64) (decimal.Decimal, error)
	CheckPayoutBlockStatus(ctx context.Context, transactionIds []int64) ([]*model.PayoutBlockStatus, error)
	CheckIncomeOperations(ctx context.Context, query model.CheckIncomeOperationsQuery) ([]*model.IncomeWithFee, error)
	CheckOthers(ctx context.Context, queryParams model.CheckOtherQuery) ([]*model.OtherOperationsWithTransaction, error)
	GetAveragePaid(ctx context.Context, queryParams model.AveragePaidQuery) (decimal.Decimal, error)
	GetServiceUserData(ctx context.Context, coinId int64, username string, limit int64) ([]*model.ServiceUserBlock, error)
}

type payouts struct {
	repo repository.PayoutRepository
}

func NewPayout(payoutRepository repository.PayoutRepository) Payouts {
	return &payouts{
		repo: payoutRepository,
	}
}

func (s *payouts) FindPayoutsForBlock(ctx context.Context, coinID int64, minPay float32, timestamp time.Time) ([]*model.PayoutForBlock, error) {

	payouts, err := s.repo.GetPayoutsForBlock(ctx, coinID, minPay, timestamp)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return payouts, nil
}

func (s *payouts) GetCurrentPayoutsBlock(ctx context.Context, coinID int64, username string, isService bool) ([]*model.PayoutBlockTransaction, error) {
	transactions, err := s.repo.GetCurrentPayoutsBlock(ctx, coinID, username, isService)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return transactions, nil
}

func (s *payouts) GetFreePayouts(ctx context.Context, coinID int64) ([]*model.FreePayout, error) {
	transactions, err := s.repo.GetFreePayouts(ctx, coinID)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return transactions, nil
}

func (s *payouts) GetCurrentPayoutsList(ctx context.Context, coinID, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error) {
	transactions, err := s.repo.GetCurrentPayoutsList(ctx, coinID, paymentTransactionType)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return transactions, nil
}

func (s *payouts) GetCurrentReferralsPayoutsList(ctx context.Context, coinID, paymentTransactionType, referralId int64) ([]*model.PayoutWithCalculation, error) {
	transactions, err := s.repo.GetCurrentReferralsPayoutsList(ctx, coinID, referralId, paymentTransactionType)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return transactions, nil
}

func (s *payouts) CheckFreePayoutTransaction(ctx context.Context, accountId, transactionId int64) (decimal.Decimal, error) {
	balance, err := s.repo.CheckFreePayoutTransaction(ctx, accountId, transactionId)

	if err != nil {
		return balance, fmt.Errorf("repository: %w", err)
	}

	return balance, nil
}

func (s *payouts) CheckPayoutBlockStatus(ctx context.Context, transactionIds []int64) ([]*model.PayoutBlockStatus, error) {
	transactions, err := s.repo.CheckPayoutBlockStatus(ctx, transactionIds)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return transactions, nil
}

func (s *payouts) CheckIncomeOperations(ctx context.Context, query model.CheckIncomeOperationsQuery) ([]*model.IncomeWithFee, error) {
	incomes, err := s.repo.CheckIncomeOperations(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return incomes, nil
}

func (s *payouts) GetAveragePaid(ctx context.Context, queryParams model.AveragePaidQuery) (decimal.Decimal, error) {
	averagePaid, err := s.repo.GetAveragePaid(ctx, queryParams)

	if err != nil {
		return averagePaid, fmt.Errorf("repository: %w", err)
	}

	return averagePaid, nil
}

func (s *payouts) CheckOthers(ctx context.Context, queryParams model.CheckOtherQuery) ([]*model.OtherOperationsWithTransaction, error) {
	operations, err := s.repo.CheckOthers(ctx, queryParams)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return operations, nil
}

func (s *payouts) GetServiceUserData(ctx context.Context, coinId int64, username string, limit int64) ([]*model.ServiceUserBlock, error) {

	blocks, err := s.repo.GetServiceUserData(ctx, coinId, username, limit)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return blocks, nil
}
