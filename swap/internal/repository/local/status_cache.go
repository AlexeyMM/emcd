package local

import (
	"sync"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
)

type statusCache struct {
	mu    sync.RWMutex
	store map[uuid.UUID]model.PublicStatus
}

func NewStatusCache() *statusCache {
	return &statusCache{
		store: make(map[uuid.UUID]model.PublicStatus),
	}
}

func (s *statusCache) Add(swapID uuid.UUID, status model.PublicStatus) {
	s.mu.Lock()
	s.store[swapID] = status
	s.mu.Unlock()
}

func (s *statusCache) Get(swapID uuid.UUID) model.PublicStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, ok := s.store[swapID]
	if !ok {
		return model.PSUnknown
	}
	return status
}

func (s *statusCache) Delete(swapID uuid.UUID) {
	s.mu.Lock()
	delete(s.store, swapID)
	s.mu.Unlock()
}
