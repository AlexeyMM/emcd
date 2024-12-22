package status_history

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type History struct {
	repository.Swap
	swapStatusHistoryRepository repository.SwapStatusHistory
}

func NewHistory(
	swapRepository repository.Swap,
	swapStatusHistoryRepository repository.SwapStatusHistory,
) *History {
	return &History{Swap: swapRepository, swapStatusHistoryRepository: swapStatusHistoryRepository}
}

func (s *History) Add(ctx context.Context, swap *model.Swap) error {
	if err := s.Swap.Add(ctx, swap); err != nil {
		return err
	}

	if err := s.swapStatusHistoryRepository.Add(ctx, swap.ID, &model.SwapStatusHistoryItem{
		Status: swap.Status, SetAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("swapStatusHistoryRepository.Add: %w", err)
	}

	return nil
}

func (s *History) Update(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial) error {
	if err := s.Swap.Update(ctx, swap, filter, partial); err != nil {
		return err
	}

	if partial.Status != nil {
		if err := s.swapStatusHistoryRepository.Add(ctx, swap.ID, &model.SwapStatusHistoryItem{
			Status: *partial.Status, SetAt: time.Now(),
		}); err != nil {
			return fmt.Errorf("swapStatusHistoryRepository.Add: %w", err)
		}
	}

	return nil
}
