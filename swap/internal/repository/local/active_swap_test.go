package local

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAddAndExist(t *testing.T) {
	activeSwap := NewActiveSwap(5)
	swapID := uuid.New()

	// Проверяем, что своп изначально не существует
	require.False(t, activeSwap.Exist(swapID), "Expected swapID %s to not exist, but it does", swapID)

	// Добавляем своп и проверяем его существование
	activeSwap.Add(swapID)
	require.True(t, activeSwap.Exist(swapID), "Expected swapID %s to exist, but it does not", swapID)
}

func TestDelete(t *testing.T) {
	activeSwap := NewActiveSwap(5)
	swapID := uuid.New()

	// Добавляем своп
	activeSwap.Add(swapID)

	// Удаляем своп
	activeSwap.Delete(swapID)

	// Проверяем, что своп был удален
	require.False(t, activeSwap.Exist(swapID), "Expected swapID %s to not exist after deletion, but it does", swapID)
}

func TestIsSwapAvailable(t *testing.T) {
	activeSwap := NewActiveSwap(2)

	// Проверяем, что своп доступен при отсутствии добавленных свопов
	require.True(t, activeSwap.IsSwapLimitExceeded(), "Expected swap to be available, but it is not")

	// Добавляем один своп
	activeSwap.Add(uuid.New())
	require.True(t, activeSwap.IsSwapLimitExceeded(), "Expected swap to be available after adding one swap, but it is not")

	// Добавляем второй своп
	activeSwap.Add(uuid.New())
	require.False(t, activeSwap.IsSwapLimitExceeded(), "Expected swap to not be available after reaching limit")
}
