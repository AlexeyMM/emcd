package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
)

type OrderFee interface {
	SyncWithAPI(ctx context.Context) error
}

type orderFee struct {
	markerCl client.Market
	rep      repository.OrderFee
}

func NewOrderFee(markerCl client.Market, rep repository.OrderFee) *orderFee {
	return &orderFee{
		markerCl: markerCl,
		rep:      rep,
	}
}

func (o *orderFee) SyncWithAPI(ctx context.Context) error {
	feeMap, err := o.markerCl.GetAllFeeRate(ctx)
	if err != nil {
		return fmt.Errorf("getAllFeeRate: %w", err)
	}
	err = o.rep.UpdateAll(ctx, feeMap)
	if err != nil {
		return fmt.Errorf("updateAll: %w", err)
	}
	return nil
}
