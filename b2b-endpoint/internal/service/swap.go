package service

import (
	"context"
	"fmt"

	internalRepository "code.emcdtech.com/b2b/endpoint/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Swap interface {
	Estimate(ctx context.Context, req *model.EstimateRequest) (amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits, err error)
	PrepareSwap(ctx context.Context, req *model.SwapRequest) (id uuid.UUID, depositAddress *model.AddressData, err error)
	StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error
	Status(ctx context.Context, swapID uuid.UUID, ch chan<- model.Status) error
	GetSwapByID(ctx context.Context, id string) (*model.SwapByIDResponse, error)
	SendSupportMessage(ctx context.Context, name, email, text string) error
}

type swap struct {
	repo      repository.Swap
	emailRepo internalRepository.Email
}

func NewSwap(repo repository.Swap, emailRepo internalRepository.Email) *swap {
	return &swap{
		repo:      repo,
		emailRepo: emailRepo,
	}
}

func (s *swap) Estimate(ctx context.Context, req *model.EstimateRequest) (amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits, err error) {
	amountFrom, amountTo, rate, limits, err = s.repo.Estimate(ctx, req)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("estimate: %w", err)
	}

	return amountFrom, amountTo, rate, limits, nil
}

func (s *swap) PrepareSwap(ctx context.Context, req *model.SwapRequest) (id uuid.UUID, depositAddress *model.AddressData, err error) {
	id, depositAddress, err = s.repo.PrepareSwap(ctx, req)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("prepareSwap: %w", err)
	}

	return id, depositAddress, nil
}

func (s *swap) StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error {
	err := s.repo.StartSwap(ctx, swapID, email, language)
	if err != nil {
		return fmt.Errorf("startSwap: %w", err)
	}
	return nil
}

func (s *swap) Status(ctx context.Context, swapID uuid.UUID, ch chan<- model.Status) error {
	err := s.repo.Status(ctx, swapID, ch)
	if err != nil {
		return fmt.Errorf("status: %v", err)
	}
	return nil
}

func (s *swap) GetSwapByID(ctx context.Context, txID string) (*model.SwapByIDResponse, error) {
	swp, err := s.repo.GetSwapByID(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("get swap by txid: %w", err)
	}

	return swp, nil
}

func (s *swap) SendSupportMessage(ctx context.Context, name, email, text string) error {
	err := s.emailRepo.SendSupportMessage(ctx, name, email, text)
	if err != nil {
		return fmt.Errorf("sendSupportMessage: %w", err)
	}
	return nil
}
