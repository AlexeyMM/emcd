package repository

import (
	"context"
	"fmt"
	"io"

	"code.emcdtech.com/b2b/swap/internal/controller/mapping"
	"code.emcdtech.com/b2b/swap/model"
	swapSwap "code.emcdtech.com/b2b/swap/protocol/swap"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/shopspring/decimal"
)

type Swap interface {
	Estimate(ctx context.Context, request *model.EstimateRequest) (amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits, err error)
	PrepareSwap(ctx context.Context, request *model.SwapRequest) (id uuid.UUID, depositAddress *model.AddressData, err error)
	StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error
	Status(ctx context.Context, swapID uuid.UUID, ch chan<- model.Status) error
	GetSwapByID(ctx context.Context, id string) (*model.SwapByIDResponse, error)
}

type swap struct {
	handler swapSwap.SwapServiceClient
}

func NewSwap(handler swapSwap.SwapServiceClient) *swap {
	return &swap{
		handler: handler,
	}
}

func (s *swap) Estimate(ctx context.Context, request *model.EstimateRequest) (amountFrom, amountTo, rate decimal.Decimal, limits *model.Limits, err error) {
	est, err := s.handler.Estimate(ctx, mapping.MapModelEstimateRequestToProto(request))
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
	}
	return mapping.MapProtoEstimateResponseToModel(est)
}

func (s *swap) PrepareSwap(ctx context.Context, request *model.SwapRequest) (id uuid.UUID, depositAddress *model.AddressData, err error) {
	response, err := s.handler.PrepareSwap(ctx, mapping.MapModelSwapRequestToProto(request))
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("swap: %w", err)
	}
	return mapping.MapProtoSwapResponseToModel(response)
}

func (s *swap) StartSwap(ctx context.Context, swapID uuid.UUID, email, language string) error {
	_, err := s.handler.StartSwap(ctx, &swapSwap.StartSwapRequest{
		SwapId:   swapID.String(),
		Email:    email,
		Language: language,
	})
	if err != nil {
		return fmt.Errorf("startSwap: %w", err)
	}
	return nil
}

func (s *swap) Status(ctx context.Context, swapID uuid.UUID, ch chan<- model.Status) error {
	stream, err := s.handler.Status(ctx, &swapSwap.StatusRequest{
		SwapId: swapID.String(),
	})
	if err != nil {
		return fmt.Errorf("status: %v", err)
	}

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case <-stream.Context().Done():
				return
			default:
				recv, err := stream.Recv()
				if err != nil && err != io.EOF {
					log.Error(ctx, "status recv: %s", err.Error())
				} else if err == io.EOF {
					return
				}
				if recv == nil {
					return
				}

				select {
				case ch <- model.Status(recv.Status):
				default:
					log.Error(ctx, "status send")
				}
			}
		}
	}()

	return nil
}

func (s *swap) GetSwapByID(ctx context.Context, id string) (*model.SwapByIDResponse, error) {
	resp, err := s.handler.GetSwapByID(ctx, &swapSwap.GetSwapByIDRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("getTransactionLink: %w", err)
	}
	return mapping.MapProtoGetSwapByIDResponseToModel(resp)
}
