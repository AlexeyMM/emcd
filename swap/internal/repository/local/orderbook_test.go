package local_test

import (
	"testing"

	businessError "code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/internal/repository/local"
	"github.com/stretchr/testify/require"
)

func TestOrderBookStore_AddSnapshot(t *testing.T) {
	store := local.NewOrderBookStore()

	symbol := "BTCUSDT"
	bids := [][2]string{
		{"50000", "1"},
		{"49900", "2"},
	}
	asks := [][2]string{
		{"51000", "1.5"},
		{"51500", "3"},
	}

	err := store.AddSnapshot(symbol, bids, asks)
	require.NoError(t, err)

	// При добавлении snapshot должны заменить старый снимок на новый
	bids = [][2]string{
		{"49700", "2"},
		{"49600", "3"},
	}
	asks = [][2]string{
		{"52000", "4"},
		{"52500", "5"},
	}

	err = store.AddSnapshot(symbol, bids, asks)
	require.NoError(t, err)

	topBids, err := store.GetBidTopLevels(symbol, 2)
	require.NoError(t, err)
	require.Len(t, topBids, 2)
	require.Equal(t, [2]float64{49700, 2}, topBids[0])
	require.Equal(t, [2]float64{49600, 3}, topBids[1])

	topAsks, err := store.GetAskTopLevels(symbol, 2)
	require.NoError(t, err)
	require.Len(t, topAsks, 2)
	require.Equal(t, [2]float64{52000, 4}, topAsks[0])
	require.Equal(t, [2]float64{52500, 5}, topAsks[1])
}

func TestOrderBookStore_AddDelta(t *testing.T) {
	store := local.NewOrderBookStore()

	symbol := "BTCUSDT"
	bids := [][2]string{
		{"50000", "1"},
		{"49900", "2"},
	}
	asks := [][2]string{
		{"51000", "1.5"},
		{"51500", "3"},
	}

	err := store.AddSnapshot(symbol, bids, asks)
	require.NoError(t, err)

	newBids := [][2]string{
		{"49900", "1.5"}, // Update size
		{"49800", "3"},   // Add new level
	}
	newAsks := [][2]string{
		{"51000", "0"}, // Remove this level
		{"52000", "2"}, // Add new level
	}

	err = store.AddDelta(symbol, newBids, newAsks)
	require.NoError(t, err)

	topBids, err := store.GetBidTopLevels(symbol, 3)
	require.NoError(t, err)
	require.Len(t, topBids, 3)
	require.Equal(t, [2]float64{50000, 1}, topBids[0])   // Unchanged
	require.Equal(t, [2]float64{49900, 1.5}, topBids[1]) // Updated
	require.Equal(t, [2]float64{49800, 3}, topBids[2])   // New level

	topAsks, err := store.GetAskTopLevels(symbol, 3)
	require.NoError(t, err)
	require.Len(t, topAsks, 2)
	require.Equal(t, [2]float64{51500, 3}, topAsks[0])
	require.Equal(t, [2]float64{52000, 2}, topAsks[1])
}

func TestOrderBookStore_GetBidTopLevels(t *testing.T) {
	store := local.NewOrderBookStore()

	symbol := "BTCUSDT"
	bids := [][2]string{
		{"50000", "1"},
		{"49900", "2"},
	}

	err := store.AddSnapshot(symbol, bids, nil)
	require.NoError(t, err)

	topBids, err := store.GetBidTopLevels(symbol, 1)
	require.NoError(t, err)
	require.Len(t, topBids, 1)
	require.Equal(t, [2]float64{50000, 1}, topBids[0])
}

func TestOrderBookStore_GetAskTopLevels(t *testing.T) {
	store := local.NewOrderBookStore()

	symbol := "BTCUSDT"
	asks := [][2]string{
		{"51000", "1.5"},
		{"51500", "3"},
	}

	err := store.AddSnapshot(symbol, nil, asks)
	require.NoError(t, err)

	topAsks, err := store.GetAskTopLevels(symbol, 1)
	require.NoError(t, err)
	require.Len(t, topAsks, 1)
	require.Equal(t, [2]float64{51000, 1.5}, topAsks[0])
}

func TestOrderBookStore_AddDelta_OrderBookNotFound(t *testing.T) {
	store := local.NewOrderBookStore()

	symbol := "BTCUSDT"
	bids := [][2]string{
		{"50000", "1"},
	}

	err := store.AddDelta(symbol, bids, nil)
	require.Error(t, err)
	require.Equal(t, businessError.OrderBookNotFoundErr, err)
}
