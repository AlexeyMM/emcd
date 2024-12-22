package service

import (
	"context"
	"fmt"
	"strings"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Symbol interface {
	SyncWithAPI(ctx context.Context) (map[string]*model.Symbol, error)
	GetAll(ctx context.Context) ([]*model.Symbol, error)
}

type symbol struct {
	repo  repository.Symbol
	byBit client.Market
}

func NewSymbol(repo repository.Symbol, byBit client.Market) *symbol {
	return &symbol{
		repo:  repo,
		byBit: byBit,
	}
}

// SyncWithAPI Получает актуальные символы по API, обновляет базу, возвращает символы
func (s *symbol) SyncWithAPI(ctx context.Context) (map[string]*model.Symbol, error) {
	symbols, err := s.byBit.GetInstrumentsInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getInstrumentsInfo: %w", err)
	}

	for _, sym := range symbols {
		sym.Accuracy = &model.Accuracy{
			BaseAccuracy:  CountDecimalPlaces(sym.BasePrecision.String()),
			QuoteAccuracy: CountDecimalPlaces(sym.QuotePrecision.String()),
		}
	}

	err = s.repo.UpdateAll(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("updateAll: %w", err)
	}

	return symbols, nil
}

func (s *symbol) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	symbols, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAll: %w", err)
	}
	return symbols, nil
}

func CountDecimalPlaces(valueStr string) int32 {
	decimalIndex := strings.Index(valueStr, ".")
	if decimalIndex == -1 {
		return 1
	}
	return int32(len(valueStr) - decimalIndex - 1)
}
