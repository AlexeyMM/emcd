package local

import (
	"context"
	"sync"

	"code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/model"
)

type Fee struct {
	mu  sync.RWMutex
	fee map[string]*model.Fee
}

func NewFee() *Fee {
	return &Fee{
		fee: make(map[string]*model.Fee),
	}
}

func (f *Fee) UpdateAll(ctx context.Context, fee map[string]*model.Fee) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.fee = fee

	return nil
}

func (f *Fee) GetFee(ctx context.Context, symbol string) (*model.Fee, error) {
	fee, ok := f.fee[symbol]
	if !ok {
		return nil, businessError.FeeNotFoundErr
	}
	return fee, nil
}
