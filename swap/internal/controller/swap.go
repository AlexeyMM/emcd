package controller

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/b2b/swap/internal/controller/mapping"
	"code.emcdtech.com/b2b/swap/internal/service"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/protocol/swap"
)

type Swap struct {
	srv        service.Swap
	subscriber service.SwapStatusSubscriber
	swap.UnimplementedSwapServiceServer
}

func NewSwap(srv service.Swap, subscriber service.SwapStatusSubscriber) *Swap {
	return &Swap{
		srv:        srv,
		subscriber: subscriber,
	}
}

func (s *Swap) Estimate(ctx context.Context, req *swap.EstimateRequest) (*swap.EstimateResponse, error) {
	estimateReq, err := mapping.MapProtoEstimateRequestToModel(req)
	if err != nil {
		log.Error(ctx, "estimate: mapping: %s", err.Error())
		return nil, fmt.Errorf("mapProtoEstimateRequestToModel: %w", err)
	}

	estimate, err := s.srv.SwapEstimate(ctx, &model.Swap{
		CoinFrom:    estimateReq.CoinFrom,
		NetworkFrom: estimateReq.NetworkFrom,
		CoinTo:      estimateReq.CoinTo,
		NetworkTo:   estimateReq.NetworkTo,
		AmountFrom:  estimateReq.AmountFrom,
		AmountTo:    estimateReq.AmountTo,
	})
	if err != nil {
		log.Error(ctx, "estimate: swapEstimate: %s", err.Error())
		return nil, fmt.Errorf("estimate: %w", err)
	}

	return mapping.MapModelEstimateResponseToProto(estimate.AmountFrom, estimate.AmountTo, estimate.Rate, estimate.Limits), nil
}

func (s *Swap) PrepareSwap(ctx context.Context, req *swap.PrepareSwapRequest) (*swap.PrepareSwapResponse, error) {
	swapRequest, err := mapping.MapProtoSwapRequestToModel(req)
	if err != nil {
		log.Error(ctx, "prepareSwap: mapping: %s", err.Error())
		return nil, fmt.Errorf("mapProtoSwapRequestToModel: %w", err)
	}

	swapID, depositAddress, err := s.srv.PrepareSwap(ctx, &model.Swap{
		CoinFrom:    swapRequest.CoinFrom,
		NetworkFrom: swapRequest.NetworkFrom,
		CoinTo:      swapRequest.CoinTo,
		AddressTo:   swapRequest.AddressTo.Address,
		NetworkTo:   swapRequest.NetworkTo,
		TagTo:       swapRequest.AddressTo.Tag,
		AmountFrom:  swapRequest.AmountFrom,
		AmountTo:    swapRequest.AmountTo,
		PartnerID:   swapRequest.PartnerID,
	})
	if err != nil {
		log.Error(ctx, "prepareSwap: %s", err.Error())
		return nil, fmt.Errorf("prepareSwap: %w", err)
	}
	return mapping.MapSwapResponseToProto(swapID, depositAddress), nil
}

func (s *Swap) StartSwap(ctx context.Context, req *swap.StartSwapRequest) (*swap.StartSwapResponse, error) {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(ctx, "startSwap: parse swap id: %s", err.Error())
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}

	err = s.srv.StartSwap(ctx, swapID, req.Email, req.Language)
	if err != nil {
		log.Error(ctx, "startSwap: %s", err.Error())
		return nil, fmt.Errorf("startSwap: %w", err)
	}
	return &swap.StartSwapResponse{}, nil
}

func (s *Swap) Status(req *swap.StatusRequest, stream swap.SwapService_StatusServer) error {
	swapID, err := uuid.Parse(req.SwapId)
	if err != nil {
		log.Error(stream.Context(), "status: parse: %s", err.Error())
		return fmt.Errorf("parse: %w", err)
	}
	clientID := uuid.New()
	ch := make(chan model.PublicStatus)

	s.subscriber.Subscribe(stream.Context(), swapID, clientID, ch)
	defer s.subscriber.Unsubscribe(swapID, clientID)

	log.Debug(stream.Context(), "status: start listen: client: %s", clientID.String())

	status, err := s.srv.GetSwapStatus(stream.Context(), swapID)
	if err != nil {
		log.Error(stream.Context(), "status: getSwapStatus: %s", err.Error())
		return fmt.Errorf("getSwapStatus: %w", err)
	}
	err = stream.Send(&swap.StatusResponse{
		Status: int32(model.ConvertInternalToPublicStatus(status)),
	})
	if err != nil {
		log.Error(stream.Context(), "status: send: %s", err.Error())
		return fmt.Errorf("send: %w", err)
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case status, ok := <-ch:
			if !ok {
				log.Debug(stream.Context(), "status: channel closed: client: %s", clientID.String())
				return nil
			}
			err = stream.Send(&swap.StatusResponse{
				Status: int32(status),
			})
			if err != nil {
				return fmt.Errorf("send: %w", err)
			}
			if status == model.PSCompleted {
				return nil
			}
		}
	}
}

func (s *Swap) GetSwapByID(ctx context.Context, req *swap.GetSwapByIDRequest) (*swap.GetSwapByIDResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Error(ctx, "getSwapByID: parse: %s", err.Error())
		return nil, fmt.Errorf("parse: %w", err)
	}

	sw, err := s.srv.GetSwapByID(ctx, id)
	if err != nil {
		log.Error(ctx, "getSwapByID: %s", err.Error())
		return nil, fmt.Errorf("getSwapByTxID: %w", err)
	}

	return mapping.MapGetSwapByIDToProto(sw), nil
}
