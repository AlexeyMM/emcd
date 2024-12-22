package local

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestErrorCounter_Inc(t *testing.T) {
	errorCounter := NewErrorCounter()
	swapID := uuid.New()

	count := errorCounter.Inc(swapID)
	require.Equal(t, 1, count)

	count = errorCounter.Inc(swapID)
	require.Equal(t, 2, count)
}

func TestErrorCounter_Delete(t *testing.T) {
	errorCounter := NewErrorCounter()
	swapID := uuid.New()

	count := errorCounter.Inc(swapID)

	count = errorCounter.Inc(swapID)
	require.Equal(t, 2, count)

	errorCounter.Delete(swapID)

	count = errorCounter.Inc(swapID)
	require.Equal(t, 1, count)
}
