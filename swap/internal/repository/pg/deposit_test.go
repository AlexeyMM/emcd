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

func TestDeposit_AddGet(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	depositRepo := NewDeposit(db)

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
	// ************************

	d1 := model.Deposit{
		TxID:      "123456",
		SwapID:    swapID,
		Coin:      "BTC",
		Amount:    decimal.NewFromFloat(123),
		Fee:       decimal.NewFromFloat(1),
		Status:    model.DepositSuccess,
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
	}
	d2 := model.Deposit{
		TxID:      "78910",
		SwapID:    swapID,
		Coin:      "ETH",
		Amount:    decimal.NewFromFloat(456),
		Fee:       decimal.NewFromFloat(2),
		Status:    model.DepositFailed,
		UpdatedAt: time.Now().Add(-24 * 2 * time.Hour).UTC(),
	} // исключаем по updatedAt

	err = depositRepo.Add(ctx, &d1)
	require.NoError(t, err)
	err = depositRepo.Add(ctx, &d2)
	require.NoError(t, err)

	tm := time.Now().Add(-24 * time.Hour)
	dps, err := depositRepo.Find(ctx, &model.DepositFilter{
		SwapID:    &swapID,
		UpdatedAt: &tm,
	})
	require.NoError(t, err)
	require.Len(t, dps, 1)
	require.Equal(t, &d1, dps[0])
}

func TestDeposit_FindOne(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	depositRepo := NewDeposit(db)

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
	// ************************

	d1 := model.Deposit{
		TxID:      "123456",
		SwapID:    swapID,
		Coin:      "BTC",
		Amount:    decimal.NewFromFloat(123),
		Fee:       decimal.NewFromFloat(1),
		Status:    model.DepositSuccess,
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
	}

	err = depositRepo.Add(ctx, &d1)
	require.NoError(t, err)

	tm := time.Now().Add(-24 * time.Hour)
	dep, err := depositRepo.FindOne(ctx, &model.DepositFilter{
		SwapID:    &swapID,
		UpdatedAt: &tm,
	})
	require.NoError(t, err)
	require.Equal(t, &d1, dep)
}
