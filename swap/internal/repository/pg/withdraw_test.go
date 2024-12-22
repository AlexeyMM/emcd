package pg

import (
	"context"
	"testing"
	"time"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestWithdraw_AddFindOne(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// References
	// ***********************
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
	}
	err := swapRepo.Add(ctx, swap)
	require.NoError(t, err)
	// ***********************

	withdrawRepo := NewWithdraw(db)

	id := int64(1)
	withdraw := &model.Withdraw{
		InternalID:         uuid.New(),
		ID:                 id,
		SwapID:             swapID,
		HashID:             "hash123",
		Coin:               "BTC",
		Network:            "mainnet",
		Address:            "1A2b3C4d5E",
		Tag:                "tag",
		Amount:             decimal.NewFromFloat(0.5),
		IncludeFeeInAmount: true,
		Status:             model.WsPending,
		ExplorerLink:       "https://explorer.example.com/tx/hash123",
	}

	err = withdrawRepo.Add(ctx, withdraw)
	require.NoError(t, err)

	retrievedWithdraw, err := withdrawRepo.FindOne(ctx, &model.WithdrawFilter{ID: &id})
	require.NoError(t, err)

	require.Equal(t, withdraw.ID, retrievedWithdraw.ID)
	require.Equal(t, withdraw.InternalID, retrievedWithdraw.InternalID)
	require.Equal(t, withdraw.SwapID, retrievedWithdraw.SwapID)
	require.Equal(t, withdraw.HashID, retrievedWithdraw.HashID)
	require.Equal(t, withdraw.Coin, retrievedWithdraw.Coin)
	require.Equal(t, withdraw.Network, retrievedWithdraw.Network)
	require.Equal(t, withdraw.Address, retrievedWithdraw.Address)
	require.Equal(t, withdraw.Tag, retrievedWithdraw.Tag)
	require.True(t, withdraw.Amount.Equal(retrievedWithdraw.Amount))
	require.Equal(t, withdraw.IncludeFeeInAmount, retrievedWithdraw.IncludeFeeInAmount)
	require.Equal(t, withdraw.Status, retrievedWithdraw.Status)
	require.Equal(t, withdraw.ExplorerLink, retrievedWithdraw.ExplorerLink)
}

func TestWithdraw_Find(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// References
	// ***********************
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
	}
	err := swapRepo.Add(ctx, swap)
	require.NoError(t, err)
	// ***********************

	withdrawRepo := NewWithdraw(db)

	withdraw := &model.Withdraw{
		ID:                 2,
		InternalID:         uuid.New(),
		SwapID:             swapID,
		HashID:             "hash456",
		Coin:               "ETH",
		Network:            "mainnet",
		Address:            "0x1234567890abcdef",
		Tag:                "tag2",
		Amount:             decimal.NewFromFloat(1.5),
		IncludeFeeInAmount: false,
		Status:             model.WsSuccess,
		ExplorerLink:       "https://explorer.example.com/tx/hash456",
	}
	withdraw2 := &model.Withdraw{
		ID:                 3,
		InternalID:         uuid.New(),
		SwapID:             swapID,
		HashID:             "323232",
		Coin:               "BTC",
		Network:            "trc20",
		Address:            "m;fddfpodfkpdfkp",
		Tag:                "tag3",
		Amount:             decimal.NewFromFloat(10.5),
		IncludeFeeInAmount: true,
		Status:             model.WsPending,
		ExplorerLink:       "https://explorer.example.com/tx/323232",
	}

	err = withdrawRepo.Add(ctx, withdraw)
	require.NoError(t, err)
	err = withdrawRepo.Add(ctx, withdraw2)
	require.NoError(t, err)

	retrievedWithdraws, err := withdrawRepo.Find(ctx, &model.WithdrawFilter{
		SwapID: &swapID,
	})
	require.NoError(t, err)

	var count int
	for _, wth := range retrievedWithdraws {
		switch wth.ID {
		case withdraw.ID:
			require.Equal(t, withdraw.ID, wth.ID)
			require.Equal(t, withdraw.InternalID, wth.InternalID)
			require.Equal(t, withdraw.SwapID, wth.SwapID)
			require.Equal(t, withdraw.HashID, wth.HashID)
			require.Equal(t, withdraw.Coin, wth.Coin)
			require.Equal(t, withdraw.Network, wth.Network)
			require.Equal(t, withdraw.Address, wth.Address)
			require.Equal(t, withdraw.Tag, wth.Tag)
			require.True(t, withdraw.Amount.Equal(wth.Amount))
			require.Equal(t, withdraw.IncludeFeeInAmount, wth.IncludeFeeInAmount)
			require.Equal(t, withdraw.Status, wth.Status)
			require.Equal(t, withdraw.ExplorerLink, wth.ExplorerLink)
			count++
		case withdraw2.ID:
			require.Equal(t, withdraw2.ID, wth.ID)
			require.Equal(t, withdraw2.InternalID, wth.InternalID)
			require.Equal(t, withdraw2.SwapID, wth.SwapID)
			require.Equal(t, withdraw2.HashID, wth.HashID)
			require.Equal(t, withdraw2.Coin, wth.Coin)
			require.Equal(t, withdraw2.Network, wth.Network)
			require.Equal(t, withdraw2.Address, wth.Address)
			require.Equal(t, withdraw2.Tag, wth.Tag)
			require.True(t, withdraw2.Amount.Equal(wth.Amount))
			require.Equal(t, withdraw2.IncludeFeeInAmount, wth.IncludeFeeInAmount)
			require.Equal(t, withdraw2.Status, wth.Status)
			require.Equal(t, withdraw2.ExplorerLink, wth.ExplorerLink)
			count++
		}
	}
	require.Equal(t, 2, count)
}

func TestWithdraw_Update(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// References
	// ***********************
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
	}
	err := swapRepo.Add(ctx, swap)
	require.NoError(t, err)
	// ***********************

	withdrawRepo := NewWithdraw(db)

	internalId := uuid.New()
	withdraw := model.Withdraw{
		ID:                 3,
		InternalID:         internalId,
		SwapID:             swapID,
		HashID:             "hash789",
		Coin:               "LTC",
		Network:            "mainnet",
		Address:            "LTCaddress",
		Tag:                "tag3",
		Amount:             decimal.NewFromFloat(2.0),
		IncludeFeeInAmount: true,
		Status:             model.WsPending,
		ExplorerLink:       "https://explorer.example.com/tx/hash789",
	}

	err = withdrawRepo.Add(ctx, &withdraw)
	require.NoError(t, err)

	newAmount := decimal.NewFromFloat(2.5)
	newStatus := model.WsSuccess
	newHashID := "newhash789"
	newID := int64(123)

	err = withdrawRepo.Update(ctx, &withdraw, &model.WithdrawFilter{
		SwapID: &swapID,
	}, &model.WithdrawPartial{
		Amount: &newAmount,
		Status: &newStatus,
		HashID: &newHashID,
		ID:     &newID,
	})
	require.NoError(t, err)

	var updatedWithdraw model.Withdraw
	err = db.QueryRow(ctx, `
		SELECT hash_id, amount, status, id
		FROM swap.withdraws
		WHERE swap_id = $1`, swapID).Scan(
		&updatedWithdraw.HashID,
		&updatedWithdraw.Amount,
		&updatedWithdraw.Status,
		&updatedWithdraw.ID,
	)
	require.NoError(t, err)
	require.Equal(t, newHashID, updatedWithdraw.HashID)
	require.True(t, newAmount.Equal(updatedWithdraw.Amount))
	require.Equal(t, newStatus, updatedWithdraw.Status)
	require.Equal(t, newID, updatedWithdraw.ID)
}
