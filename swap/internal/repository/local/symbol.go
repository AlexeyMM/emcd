package local

import (
	"context"
	"sync"

	"code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/model"
)

type Symbol struct {
	mu      sync.RWMutex
	storage map[string]*model.Symbol
}

func NewSymbol(storage map[string]*model.Symbol) *Symbol {
	return &Symbol{
		storage: storage,
	}
}

func (s *Symbol) UpdateAll(ctx context.Context, symbols map[string]*model.Symbol) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage = symbols

	return nil
}

func (s *Symbol) Get(ctx context.Context, title string) (*model.Symbol, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	symbol, ok := s.storage[title]
	if !ok {
		return nil, businessError.SymbolNotFoundErr
	}
	return symbol, nil
}

func (s *Symbol) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	allSymbols := make([]*model.Symbol, 0, len(s.storage))
	for _, symbol := range s.storage {
		allSymbols = append(allSymbols, symbol)
	}
	return allSymbols, nil
}

func (s *Symbol) GetAccuracy(ctx context.Context, symbol string) (*model.Accuracy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sym, ok := s.storage[symbol]
	if !ok {
		return &model.Accuracy{}, nil
	}
	return sym.Accuracy, nil
}
