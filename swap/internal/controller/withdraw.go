package controller

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/protocol/swapWithdraw"
	"github.com/google/uuid"
)

type Withdraw struct {
	srv service.Withdraw
	swapWithdraw.UnimplementedSwapWithdrawServiceServer
}

func NewWithdraw(srv service.Withdraw) *Withdraw {
	return &Withdraw{
		srv: srv,
	}
}

func (w *Withdraw) GetTransactionLink(ctx context.Context, req *swapWithdraw.GetTransactionLinkRequest) (*swapWithdraw.GetTransactionLinkResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	sw, err := w.srv.GetBySwapID(ctx, swapID)
	if err != nil {
		return nil, fmt.Errorf("getBySwapID: %w", err)
	}

	return &swapWithdraw.GetTransactionLinkResponse{
		Link: sw.ExplorerLink,
	}, nil
}
