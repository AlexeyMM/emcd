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

func TestOrder_AddGetUpdate(t *testing.T) {
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

	accountID := int64(120)
	acc := model.Account{
		ID: int64(accountID),
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}
	rep := NewAccount(db)
	err = rep.Add(ctx, &acc)
	require.NoError(t, err)
	// ************************

	orderRepo := NewOrder(db)

	order := model.Order{
		ID:         uuid.New(),
		SwapID:     swapID,
		AccountID:  accountID,
		Category:   "spot",
		Symbol:     "BTCUSD",
		Direction:  model.Buy,
		AmountFrom: decimal.NewFromFloat(0.5),
		AmountTo:   decimal.NewFromFloat(25000),
		Status:     model.OrderPending,
		IsFirst:    true,
	}

	err = orderRepo.Add(ctx, &order)
	require.NoError(t, err)

	retrievedOrder, err := orderRepo.FindOne(ctx, &model.OrderFilter{
		ID: &order.ID,
	})
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, order.ID, retrievedOrder.ID)
	require.Equal(t, order.SwapID, retrievedOrder.SwapID)
	require.Equal(t, order.AccountID, retrievedOrder.AccountID)
	require.Equal(t, order.Category, retrievedOrder.Category)
	require.Equal(t, order.Symbol, retrievedOrder.Symbol)
	require.Equal(t, order.Direction, retrievedOrder.Direction)
	require.True(t, order.AmountFrom.Equal(retrievedOrder.AmountFrom))
	require.True(t, order.AmountTo.Equal(retrievedOrder.AmountTo))
	require.Equal(t, order.Status, retrievedOrder.Status)
	require.Equal(t, order.IsFirst, retrievedOrder.IsFirst)

	newStatus := model.OrderFilled
	newAmountFrom := decimal.NewFromFloat(0.49)
	newAmountTo := decimal.NewFromFloat(25001)
	err = orderRepo.Update(ctx, &order,
		&model.OrderFilter{
			ID: &order.ID,
		}, &model.OrderPartial{
			AmountFrom: &newAmountFrom,
			AmountTo:   &newAmountTo,
			Status:     &newStatus,
		})
	require.NoError(t, err)

	updatedOrder, err := orderRepo.Find(ctx, &model.OrderFilter{
		ID:        &order.ID,
		AccountID: &order.AccountID,
		IsFirst:   &order.IsFirst,
	})
	require.NoError(t, err)

	require.Equal(t, &order, updatedOrder[0])
}

func TestOrder_UpdateStatus(t *testing.T) {
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

	accountID := int64(121)
	acc := model.Account{
		ID: int64(accountID),
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}
	rep := NewAccount(db)
	err = rep.Add(ctx, &acc)
	require.NoError(t, err)
	// ************************

	orderRepo := NewOrder(db)

	order := model.Order{
		ID:         uuid.New(),
		SwapID:     swapID,
		AccountID:  accountID,
		Category:   "spot",
		Symbol:     "BTCUSD",
		Direction:  model.Buy,
		AmountFrom: decimal.NewFromFloat(0.5),
		AmountTo:   decimal.NewFromFloat(25000),
		Status:     model.OrderCreated,
		IsFirst:    true,
	}

	err = orderRepo.Add(ctx, &order)
	require.NoError(t, err)

	newStatus := model.OrderFilled
	err = orderRepo.Update(ctx, &order,
		&model.OrderFilter{
			ID: &order.ID,
		}, &model.OrderPartial{
			Status: &newStatus,
		})
	require.NoError(t, err)

	updatedOrder, err := orderRepo.FindOne(ctx, &model.OrderFilter{
		ID: &order.ID,
	})
	require.NoError(t, err)

	require.Equal(t, newStatus, updatedOrder.Status)
}

func TestOrder_UpdateStatusNotUpdate(t *testing.T) {
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

	accountID := int64(122)
	acc := model.Account{
		ID: int64(accountID),
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}
	rep := NewAccount(db)
	err = rep.Add(ctx, &acc)
	require.NoError(t, err)
	// ************************

	orderRepo := NewOrder(db)

	order := model.Order{
		ID:         uuid.New(),
		SwapID:     swapID,
		AccountID:  accountID,
		Category:   "spot",
		Symbol:     "BTCUSD",
		Direction:  model.Buy,
		AmountFrom: decimal.NewFromFloat(0.5),
		AmountTo:   decimal.NewFromFloat(25000),
		Status:     model.OrderFilled,
		IsFirst:    true,
	}

	err = orderRepo.Add(ctx, &order)
	require.NoError(t, err)

	newStatus := model.OrderCreated
	err = orderRepo.Update(ctx, &order,
		&model.OrderFilter{
			ID: &order.ID,
		}, &model.OrderPartial{
			Status: &newStatus,
		})
	require.NoError(t, err)

	updatedOrder, err := orderRepo.FindOne(ctx,
		&model.OrderFilter{
			ID: &order.ID,
		})
	require.NoError(t, err)

	require.Equal(t, newStatus, updatedOrder.Status)
}

func TestOrder_Find(t *testing.T) {
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

	accountID := int64(123)
	acc := model.Account{
		ID: int64(accountID),
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}
	rep := NewAccount(db)
	err = rep.Add(ctx, &acc)
	require.NoError(t, err)
	// ************************

	orderRepo := NewOrder(db)

	order := model.Order{
		ID:         uuid.New(),
		SwapID:     swapID,
		AccountID:  accountID,
		Category:   "spot",
		Symbol:     "BTCUSD",
		Direction:  model.Buy,
		AmountFrom: decimal.NewFromFloat(0.5),
		AmountTo:   decimal.NewFromFloat(25000),
		Status:     model.OrderPending,
		IsFirst:    true,
	}

	err = orderRepo.Add(ctx, &order)
	require.NoError(t, err)

	order2 := model.Order{
		ID:         uuid.New(),
		SwapID:     swapID,
		AccountID:  accountID,
		Category:   "spot",
		Symbol:     "BTCUSD",
		Direction:  model.Sell,
		AmountFrom: decimal.NewFromFloat(0.5),
		AmountTo:   decimal.NewFromFloat(25000),
		Status:     model.OrderPending,
		IsFirst:    false,
	}

	err = orderRepo.Add(ctx, &order2)
	require.NoError(t, err)

	isFirst := true

	retrievedOrder, err := orderRepo.FindOne(ctx,
		&model.OrderFilter{
			AccountID: &accountID,
			IsFirst:   &isFirst,
		})
	require.NoError(t, err)

	require.Equal(t, order.ID, retrievedOrder.ID)
	require.Equal(t, order.SwapID, retrievedOrder.SwapID)
	require.Equal(t, order.AccountID, retrievedOrder.AccountID)
	require.Equal(t, order.Category, retrievedOrder.Category)
	require.Equal(t, order.Symbol, retrievedOrder.Symbol)
	require.Equal(t, order.Direction, retrievedOrder.Direction)
	require.True(t, order.AmountFrom.Equal(retrievedOrder.AmountFrom))
	require.True(t, order.AmountTo.Equal(retrievedOrder.AmountTo))
	require.Equal(t, order.Status, retrievedOrder.Status)
	require.Equal(t, order.IsFirst, retrievedOrder.IsFirst)

	isFirst = false
	retrievedOrder, err = orderRepo.FindOne(ctx,
		&model.OrderFilter{
			AccountID: &accountID,
			IsFirst:   &isFirst,
		})
	require.NoError(t, err)

	require.Equal(t, order2.ID, retrievedOrder.ID)
	require.Equal(t, order2.SwapID, retrievedOrder.SwapID)
	require.Equal(t, order2.AccountID, retrievedOrder.AccountID)
	require.Equal(t, order2.Category, retrievedOrder.Category)
	require.Equal(t, order2.Symbol, retrievedOrder.Symbol)
	require.Equal(t, order2.Direction, retrievedOrder.Direction)
	require.True(t, order2.AmountFrom.Equal(retrievedOrder.AmountFrom))
	require.True(t, order2.AmountTo.Equal(retrievedOrder.AmountTo))
	require.Equal(t, order2.Status, retrievedOrder.Status)
	require.Equal(t, order2.IsFirst, retrievedOrder.IsFirst)
}
