package local

import (
	"sync"

	"github.com/google/uuid"
)

type ErrorCounter struct {
	mu    sync.Mutex
	swaps map[uuid.UUID]int
}

func NewErrorCounter() *ErrorCounter {
	return &ErrorCounter{
		swaps: make(map[uuid.UUID]int),
	}
}

func (e *ErrorCounter) Inc(swapID uuid.UUID) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.swaps[swapID]++

	return e.swaps[swapID]
}

func (e *ErrorCounter) Delete(swapID uuid.UUID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.swaps, swapID)
}
