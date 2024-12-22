package repository

import "github.com/google/uuid"

// ActiveSwap за Add, Delete отвечает worker.swap_executor
type ActiveSwap interface {
	Add(swapID uuid.UUID)
	Exist(swapID uuid.UUID) bool
	Delete(swapID uuid.UUID)
	// IsSwapLimitExceeded Своп доступен для выполнения.
	// Работает на основании максимально допустимого количества одновременно выполняемых свопов.
	IsSwapLimitExceeded() bool
}
