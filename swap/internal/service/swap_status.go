package service

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type SwapStatusUpdater interface {
	GetSwap(ctx context.Context, swapID uuid.UUID) (*model.Swap, error)
	UpdateAndBroadcast(ctx context.Context, swap *model.Swap, status model.Status) error
	Update(ctx context.Context, swap *model.Swap, status model.Status) error
	Broadcast(ctx context.Context, swap *model.Swap, status model.Status)
}

type SwapStatusSubscriber interface {
	Subscribe(ctx context.Context, swapID, clientID uuid.UUID, ch chan model.PublicStatus)
	Unsubscribe(swapID, clientID uuid.UUID)
}

type SwapStatus struct {
	swap              repository.Swap
	swapStatusHistory repository.SwapStatusHistory
	subscribers       repository.Subscribers
	statusCache       repository.StatusCache
}

func NewSwapStatus(
	swap repository.Swap,
	swapStatusHistory repository.SwapStatusHistory,
	subscribers repository.Subscribers,
	statusCache repository.StatusCache,
) *SwapStatus {
	return &SwapStatus{
		swap:              swap,
		swapStatusHistory: swapStatusHistory,
		subscribers:       subscribers,
		statusCache:       statusCache,
	}
}

func (s *SwapStatus) GetSwap(ctx context.Context, swapID uuid.UUID) (*model.Swap, error) {
	sw, err := s.swap.FindOne(ctx, &model.SwapFilter{
		ID: &swapID,
	})
	if err != nil {
		return nil, fmt.Errorf("findOne: %w", err)
	}
	return sw, nil
}

func (s *SwapStatus) Update(ctx context.Context, swap *model.Swap, status model.Status) error {
	if err := s.swap.Update(ctx, swap,
		&model.SwapFilter{
			ID: &swap.ID,
		},
		getSwapPartial(status)); err != nil {
		return fmt.Errorf("swapRep.Update: %w", err)
	}

	return nil
}

// UpdateAndBroadcast обновляет базу, кеш, отправляет статус клиенту
func (s *SwapStatus) UpdateAndBroadcast(ctx context.Context, swap *model.Swap, status model.Status) error {
	if err := s.Update(ctx, swap, status); err != nil {
		return fmt.Errorf("swapStatus.Update: %w", err)
	}

	lastSentStatus := s.statusCache.Get(swap.ID)
	if lastSentStatus != model.ConvertInternalToPublicStatus(status) {
		s.broadcast(ctx, swap.ID, model.ConvertInternalToPublicStatus(status))
		s.statusCache.Add(swap.ID, model.ConvertInternalToPublicStatus(status))
	}

	return nil
}

// Broadcast не обновляет базу, обновляет кеш, отправляет статус клиенту.
// Используется когда нужно показать Error клиенту, но мы не хотим сетить error в базу,
// потому что рассчитываем пофиксить баг, перезапустить сервис и начать с того же шага.
func (s *SwapStatus) Broadcast(ctx context.Context, swap *model.Swap, status model.Status) {
	if s.statusCache.Get(swap.ID) != model.ConvertInternalToPublicStatus(status) {
		s.broadcast(ctx, swap.ID, model.ConvertInternalToPublicStatus(status))
		s.statusCache.Add(swap.ID, model.ConvertInternalToPublicStatus(status))
	}
}

func (s *SwapStatus) Subscribe(ctx context.Context, swapID, clientID uuid.UUID, ch chan model.PublicStatus) {
	s.subscribers.Add(swapID, clientID, ch)
}

func (s *SwapStatus) Unsubscribe(swapID, clientID uuid.UUID) {
	s.subscribers.Delete(swapID, clientID)

	subs := s.subscribers.GetBySwapID(swapID)
	if subs == nil {
		s.statusCache.Delete(swapID)
	}
}

func (s *SwapStatus) broadcast(ctx context.Context, swapID uuid.UUID, status model.PublicStatus) {
	subs := s.subscribers.GetBySwapID(swapID)
	if subs == nil {
		return
	}

	for _, sub := range subs {
		select {
		case sub.Ch <- status:
		default:
			log.Error(ctx, "failed to broadcast status to client: swapID %s client: %s",
				swapID.String(), sub.ClientID.String())
		}
	}
}

func getSwapPartial(status model.Status) *model.SwapPartial {
	partial := model.SwapPartial{
		Status: &status,
	}

	switch status {
	case model.Completed, model.Cancel, model.Error, model.ManualCompleted:
		endTime := time.Now().UTC()
		partial.EndTime = &endTime
	}

	return &partial
}
