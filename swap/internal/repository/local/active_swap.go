package local

import (
	"sync"

	"github.com/google/uuid"
)

type ActiveSwap struct {
	mu      sync.RWMutex
	swapIDs map[uuid.UUID]struct{}
	limit   int // Максимально доступное количество свопов
}

func NewActiveSwap(limit int) *ActiveSwap {
	return &ActiveSwap{
		swapIDs: make(map[uuid.UUID]struct{}),
		limit:   limit,
	}
}

func (a *ActiveSwap) Add(swapID uuid.UUID) {
	a.mu.Lock()
	a.swapIDs[swapID] = struct{}{}
	a.mu.Unlock()
}

func (a *ActiveSwap) Exist(swapID uuid.UUID) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, ok := a.swapIDs[swapID]
	return ok
}

func (a *ActiveSwap) Delete(swapID uuid.UUID) {
	a.mu.Lock()
	delete(a.swapIDs, swapID)
	a.mu.Unlock()
}

func (a *ActiveSwap) IsSwapLimitExceeded() bool {
	return len(a.swapIDs) < a.limit
}
