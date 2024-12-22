package pg

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/b2b/swap/model"
)

func TestSwapStatusHistory_Add(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	swapId := uuid.New()

	swapRepo := NewSwap(db)

	swapStatusHistoryRepo := NewSwapStatusHistory(db)

	historyItem := &model.SwapStatusHistoryItem{
		Status: model.WaitDeposit, SetAt: time.Now().Truncate(time.Second),
	}
	err := swapStatusHistoryRepo.Add(ctx, swapId, historyItem)
	require.Error(t, err, "adding swap history for swap which doesn't exist should return error")

	err = swapRepo.Add(ctx, &model.Swap{ID: swapId})
	require.NoError(t, err)

	err = swapStatusHistoryRepo.Add(ctx, swapId, historyItem)
	require.NoError(t, err, "adding swap history for swap which exists should not return error")

	find, err := swapStatusHistoryRepo.Find(ctx, &model.SwapStatusHistoryFilter{
		SwapID: &swapId,
	})
	require.NoError(t, err)
	require.Len(t, find, 1)
	require.Equal(t, historyItem.Status, find[0].Status)
	require.Equal(t, historyItem.SetAt.UTC(), find[0].SetAt)

	secondHistoryItem := &model.SwapStatusHistoryItem{
		Status: model.CheckDeposit,
		SetAt:  time.Now().Add(time.Second).Truncate(time.Second),
	}

	err = swapStatusHistoryRepo.Add(ctx, swapId, secondHistoryItem)
	require.NoError(t, err, "adding swap history for swap which exists should not return error")

	find, err = swapStatusHistoryRepo.Find(ctx, &model.SwapStatusHistoryFilter{
		SwapID: &swapId,
	})
	require.NoError(t, err)
	require.Len(t, find, 2)
	require.Equal(t, secondHistoryItem.Status, find[0].Status)
	require.Equal(t, secondHistoryItem.SetAt.UTC(), find[0].SetAt)
	require.Equal(t, historyItem.Status, find[1].Status)
	require.Equal(t, historyItem.SetAt.UTC(), find[1].SetAt)
}

func TestSwapStatusHistory_Find(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	swapId := uuid.New()
	secondSwapId := uuid.New()

	swapRepo := NewSwap(db)
	swapStatusHistoryRepo := NewSwapStatusHistory(db)

	historyItem := &model.SwapStatusHistoryItem{
		Status: model.WaitDeposit, SetAt: time.Now().Truncate(time.Second),
	}

	err := swapRepo.Add(ctx, &model.Swap{ID: swapId})
	require.NoError(t, err)

	err = swapRepo.Add(ctx, &model.Swap{ID: secondSwapId})
	require.NoError(t, err)

	err = swapStatusHistoryRepo.Add(ctx, swapId, historyItem)
	require.NoError(t, err, "adding swap history for swap which exists should not return error")

	err = swapStatusHistoryRepo.Add(ctx, secondSwapId, historyItem)
	require.NoError(t, err, "adding swap history for swap which exists should not return error")

	find, err := swapStatusHistoryRepo.Find(ctx, &model.SwapStatusHistoryFilter{
		SwapID: &swapId,
	})
	require.NoError(t, err)
	require.Len(t, find, 1)
	require.Equal(t, historyItem.Status, find[0].Status)
	require.Equal(t, historyItem.SetAt.UTC(), find[0].SetAt)
}
