package pg

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/b2b/swap/model"
)

func TestSwap_Find(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	swapRepo := NewSwap(db)

	swaps := []*model.Swap{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 1,
			CoinFrom:    "BTC",
			AddressFrom: "1A2B3C4D5E6F7G8H9I0J",
			NetworkFrom: "mainnet1",
			CoinTo:      "ETH",
			AddressTo:   "0xabc123...",
			NetworkTo:   "mainnet2",
			TagTo:       "tag123",
			AmountFrom:  decimal.NewFromFloat(1.5),
			AmountTo:    decimal.NewFromFloat(2500),
			Status:      model.WaitDeposit,
			StartTime:   time.Now().UTC().Truncate(time.Second),
			EndTime:     time.Now().UTC().Add(time.Hour).Truncate(time.Second),
			PartnerID:   "partner 1",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 2,
			CoinFrom:    "ETH",
			AddressFrom: "0xdef456...",
			NetworkFrom: "mainnet2",
			CoinTo:      "BTC",
			AddressTo:   "1J2K3L4M5N6O7P8Q9R0S",
			NetworkTo:   "mainnet1",
			TagTo:       "tag456",
			AmountFrom:  decimal.NewFromFloat(2.0),
			AmountTo:    decimal.NewFromFloat(3000),
			Status:      model.Completed, // должен быть исключен
			StartTime:   time.Now().UTC(),
			EndTime:     time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second),
			PartnerID:   "partner 2",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 3,
			CoinFrom:    "LTC",
			AddressFrom: "1X2Y3Z4W5V6U7T8S9R0Q",
			NetworkFrom: "mainnet3",
			CoinTo:      "BTC",
			AddressTo:   "1M2N3O4P5Q6R7S8T9U0V",
			NetworkTo:   "mainnet1",
			TagTo:       "tag789",
			AmountFrom:  decimal.NewFromFloat(3.5),
			AmountTo:    decimal.NewFromFloat(4000),
			Status:      model.CheckOrder,
			StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(time.Second),
			EndTime:     time.Now().UTC().Add(3 * time.Hour).Truncate(time.Second),
			PartnerID:   "partner 3",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 4,
			CoinFrom:    "BTC",
			AddressFrom: "DFEFEFEFEFEFFE",
			NetworkFrom: "DFFDSFD",
			CoinTo:      "BTC",
			AddressTo:   "DFP,DFDF,[PDF,DF",
			NetworkTo:   "mainnet4",
			TagTo:       "tag78933434",
			AmountFrom:  decimal.NewFromFloat(4.5),
			AmountTo:    decimal.NewFromFloat(4001),
			Status:      model.Cancel,
			StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(2 * time.Second),
			EndTime:     time.Now().UTC().Add(4 * time.Hour).Truncate(2 * time.Second),
			PartnerID:   "partner 4",
		},
	}

	for _, swap := range swaps {
		err := swapRepo.Add(ctx, swap)
		require.NoError(t, err)
	}

	statusCompleted := []model.Status{model.Completed, model.Cancel}
	activeSwaps, err := swapRepo.Find(ctx, &model.SwapFilter{
		NotEqStatus: statusCompleted,
	})
	require.NoError(t, err)
	require.Len(t, activeSwaps, 2)

	for _, swap := range activeSwaps {
		require.NotEqual(t, model.Completed, swap.Status)
		found := false
		for _, expectedSwap := range swaps {
			if swap.ID == expectedSwap.ID {
				require.Equal(t, expectedSwap.AccountFrom, swap.AccountFrom)
				require.Equal(t, expectedSwap.CoinFrom, swap.CoinFrom)
				require.Equal(t, expectedSwap.AddressFrom, swap.AddressFrom)
				require.Equal(t, expectedSwap.NetworkFrom, swap.NetworkFrom)
				require.Equal(t, expectedSwap.CoinTo, swap.CoinTo)
				require.Equal(t, expectedSwap.AddressTo, swap.AddressTo)
				require.Equal(t, expectedSwap.NetworkTo, swap.NetworkTo)
				require.Equal(t, expectedSwap.TagTo, swap.TagTo)
				require.True(t, expectedSwap.AmountFrom.Equal(swap.AmountFrom))
				require.True(t, expectedSwap.AmountTo.Equal(swap.AmountTo))
				require.Equal(t, expectedSwap.Status, swap.Status)
				require.Equal(t, expectedSwap.StartTime, swap.StartTime)
				require.Equal(t, expectedSwap.EndTime, swap.EndTime)
				require.Equal(t, expectedSwap.PartnerID, swap.PartnerID)
				found = true
				break
			}
		}
		require.True(t, found, "Swap not found in the expected list")
	}

}

func TestSwap_FindOne(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	swapRepo := NewSwap(db)

	swapID := uuid.New()
	swap := &model.Swap{
		ID:          swapID,
		AccountFrom: 1,
		CoinFrom:    "BTC",
		AddressFrom: "1A2B3C4D5E6F7G8H9I0J",
		NetworkFrom: "mainnet1",
		CoinTo:      "ETH",
		AddressTo:   "0xabc123...",
		NetworkTo:   "mainnet2",
		TagTo:       "tag123",
		AmountFrom:  decimal.NewFromFloat(1.5),
		AmountTo:    decimal.NewFromFloat(2500),
		Status:      model.WaitDeposit,
		StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(time.Second),
		EndTime:     time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second),
		PartnerID:   "partner 1",
	}

	err := swapRepo.Add(ctx, swap)
	require.NoError(t, err)

	retrievedSwap, err := swapRepo.FindOne(ctx, &model.SwapFilter{
		ID: &swapID,
	})
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, swap.ID, retrievedSwap.ID)
	require.Equal(t, swap.AccountFrom, retrievedSwap.AccountFrom)
	require.Equal(t, swap.CoinFrom, retrievedSwap.CoinFrom)
	require.Equal(t, swap.AddressFrom, retrievedSwap.AddressFrom)
	require.Equal(t, swap.NetworkFrom, retrievedSwap.NetworkFrom)
	require.Equal(t, swap.CoinTo, retrievedSwap.CoinTo)
	require.Equal(t, swap.AddressTo, retrievedSwap.AddressTo)
	require.Equal(t, swap.NetworkTo, retrievedSwap.NetworkTo)
	require.Equal(t, swap.TagTo, retrievedSwap.TagTo)
	require.True(t, swap.AmountFrom.Equal(retrievedSwap.AmountFrom))
	require.True(t, swap.AmountTo.Equal(retrievedSwap.AmountTo))
	require.Equal(t, swap.Status, retrievedSwap.Status)
	require.Equal(t, swap.StartTime, retrievedSwap.StartTime)
	require.Equal(t, swap.EndTime, retrievedSwap.EndTime)
	require.Equal(t, swap.PartnerID, retrievedSwap.PartnerID)
}

func TestSwap_Update(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	swapRepo := NewSwap(db)

	swapID := uuid.New()
	swap := model.Swap{
		ID:          swapID,
		UserID:      uuid.Nil,
		AccountFrom: 1,
		CoinFrom:    "BTC",
		AddressFrom: "1A2B3C4D5E6F7G8H9I0J",
		TagFrom:     "12",
		NetworkFrom: "mainnet1",
		CoinTo:      "ETH",
		AddressTo:   "0xabc123...",
		NetworkTo:   "mainnet2",
		TagTo:       "tag123",
		AmountFrom:  decimal.NewFromFloat(1.5),
		AmountTo:    decimal.NewFromFloat(25),
		Status:      model.WaitDeposit,
		StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(time.Second),
		EndTime:     time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second),
	}

	err := swapRepo.Add(ctx, &swap)
	require.NoError(t, err)

	userID := uuid.New()
	statusCompleted := model.Completed
	startTime := time.Now().UTC().Truncate(time.Second)
	endTime := time.Now().UTC().Add(5 * time.Hour).Truncate(time.Second)
	amountFrom := decimal.NewFromFloat(2.5)
	amountTo := decimal.NewFromFloat(35)
	addressTo := "new address to"
	tagTo := "new tag to"

	err = swapRepo.Update(ctx, &swap,
		&model.SwapFilter{
			ID: &swapID,
		}, &model.SwapPartial{
			UserID:     &userID,
			Status:     &statusCompleted,
			AmountFrom: &amountFrom,
			AmountTo:   &amountTo,
			StartTime:  &startTime,
			EndTime:    &endTime,
			AddressTo:  &addressTo,
			TagTo:      &tagTo,
		})
	require.NoError(t, err)

	var (
		updatedUserID     uuid.UUID
		status            model.Status
		updatedStartTime  time.Time
		updatedEndTime    time.Time
		updatedAmountFrom decimal.Decimal
		updatedAmountTo   decimal.Decimal
		updatedAddressTo  string
		updatedTagTo      string
	)
	err = db.QueryRow(ctx, `
		SELECT user_id, status, start_time, end_time, amount_from, amount_to, address_to, tag_to
		FROM swap.swaps
		WHERE id = $1`, swapID).Scan(
		&updatedUserID,
		&status,
		&updatedStartTime,
		&updatedEndTime,
		&updatedAmountFrom,
		&updatedAmountTo,
		&updatedAddressTo,
		&updatedTagTo,
	)
	require.NoError(t, err)
	require.Equal(t, userID, updatedUserID)
	require.Equal(t, statusCompleted, status)
	require.Equal(t, startTime, updatedStartTime)
	require.Equal(t, endTime, updatedEndTime)
	require.Equal(t, amountFrom, updatedAmountFrom)
	require.Equal(t, amountTo, updatedAmountTo)
	require.Equal(t, addressTo, updatedAddressTo)
	require.Equal(t, tagTo, updatedTagTo)
}

func TestSwap_GetCountSwapsByStatus(t *testing.T) {
	ctx := context.Background()

	swapRepo := NewSwap(db)

	swaps := []*model.Swap{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 1,
			CoinFrom:    "BTC",
			AddressFrom: "1A2B3C4D5E6F7G8H9I0J",
			NetworkFrom: "mainnet1",
			CoinTo:      "ETH",
			AddressTo:   "0xabc123...",
			NetworkTo:   "mainnet2",
			TagTo:       "tag123",
			AmountFrom:  decimal.NewFromFloat(1.5),
			AmountTo:    decimal.NewFromFloat(2500),
			Status:      model.WaitDeposit,
			StartTime:   time.Now().UTC().Truncate(time.Second),
			EndTime:     time.Now().UTC().Add(time.Hour).Truncate(time.Second),
			PartnerID:   "partner 1",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 2,
			CoinFrom:    "ETH",
			AddressFrom: "0xdef456...",
			NetworkFrom: "mainnet2",
			CoinTo:      "BTC",
			AddressTo:   "1J2K3L4M5N6O7P8Q9R0S",
			NetworkTo:   "mainnet1",
			TagTo:       "tag456",
			AmountFrom:  decimal.NewFromFloat(2.0),
			AmountTo:    decimal.NewFromFloat(3000),
			Status:      model.WaitDeposit,
			StartTime:   time.Now().UTC(),
			EndTime:     time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second),
			PartnerID:   "partner 2",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 3,
			CoinFrom:    "LTC",
			AddressFrom: "1X2Y3Z4W5V6U7T8S9R0Q",
			NetworkFrom: "mainnet3",
			CoinTo:      "BTC",
			AddressTo:   "1M2N3O4P5Q6R7S8T9U0V",
			NetworkTo:   "mainnet1",
			TagTo:       "tag789",
			AmountFrom:  decimal.NewFromFloat(3.5),
			AmountTo:    decimal.NewFromFloat(4000),
			Status:      model.WaitWithdraw,
			StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(time.Second),
			EndTime:     time.Now().UTC().Add(3 * time.Hour).Truncate(time.Second),
			PartnerID:   "partner 3",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			AccountFrom: 4,
			CoinFrom:    "BTC",
			AddressFrom: "DFEFEFEFEFEFFE",
			NetworkFrom: "DFFDSFD",
			CoinTo:      "BTC",
			AddressTo:   "DFP,DFDF,[PDF,DF",
			NetworkTo:   "mainnet4",
			TagTo:       "tag78933434",
			AmountFrom:  decimal.NewFromFloat(4.5),
			AmountTo:    decimal.NewFromFloat(4001),
			Status:      model.Completed,
			StartTime:   time.Now().UTC().Add(-time.Hour).Truncate(2 * time.Second),
			EndTime:     time.Now().UTC().Add(4 * time.Hour).Truncate(2 * time.Second),
			PartnerID:   "partner 4",
		},
	}

	for _, swap := range swaps {
		err := swapRepo.Add(ctx, swap)
		require.NoError(t, err)
	}

	statuses, err := swapRepo.CountSwapsByStatus(ctx)
	require.NoError(t, err)

	expectedStatuses := map[model.Status]int{
		model.Completed:    1,
		model.WaitWithdraw: 1,
		model.WaitDeposit:  2,
	}

	require.Equal(t, expectedStatuses, statuses)
}
