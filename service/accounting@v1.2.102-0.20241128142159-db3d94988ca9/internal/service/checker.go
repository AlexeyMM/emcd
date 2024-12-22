package service

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"fmt"
	"time"
)

type Checker interface {
	GetIncomesHashrateByDate(ctx context.Context, date time.Time) ([]*model.HashrateByDate, error)
	GetCoinsOperationsSum(ctx context.Context) ([]*model.SumCheckData, error)
	GetTransactionOperationsIntegrity(ctx context.Context) ([]*model.TransactionOperationsIntegrityData, error)
	GetCheckTransactionCoins(ctx context.Context) (model.CheckTransactionCoinsData, error)
	GetCheckFreezePayoutsBlocks(ctx context.Context) ([]*model.CheckFreezePayoutsBlocksData, error)
}

type checker struct {
	repo repository.CheckerRepository
}

func NewChecker(checkerRepository repository.CheckerRepository) Checker {
	return &checker{
		repo: checkerRepository,
	}
}

func (s *checker) GetIncomesHashrateByDate(ctx context.Context, date time.Time) ([]*model.HashrateByDate, error) {

	hashrateData, err := s.repo.GetIncomesHashrateByDate(ctx, date)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return hashrateData, nil
}

func (s *checker) GetCoinsOperationsSum(ctx context.Context) ([]*model.SumCheckData, error) {

	result, err := s.repo.GetCoinsOperationsSum(ctx)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return result, nil
}

func (s *checker) GetTransactionOperationsIntegrity(ctx context.Context) ([]*model.TransactionOperationsIntegrityData, error) {
	result, err := s.repo.GetTransactionOperationsIntegrity(ctx)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return result, nil
}

func (s *checker) GetCheckTransactionCoins(ctx context.Context) (model.CheckTransactionCoinsData, error) {

	result := model.CheckTransactionCoinsData{
		TrIds: nil,
		OpIds: nil,
	}

	resultTr, err := s.repo.GetCheckTransactionCoins(ctx)

	if err != nil {
		return result, fmt.Errorf("repository: %w", err)
	}

	if len(resultTr) > 0 {
		result.TrIds = resultTr
	}

	resultOp, err := s.repo.GetCheckOperationsCoins(ctx)

	if err != nil {
		return result, fmt.Errorf("repository: %w", err)
	}

	if len(resultOp) > 0 {
		result.OpIds = resultOp
	}

	return result, nil
}

func (s *checker) GetCheckFreezePayoutsBlocks(ctx context.Context) ([]*model.CheckFreezePayoutsBlocksData, error) {

	result, err := s.repo.GetCheckFreezePayoutsBlocks(ctx)

	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	return result, nil
}
