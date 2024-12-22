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

func TestTransfer_AddFind(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	transferRepo := NewTransfer(db)

	fromAccountID := int64(10)
	transfer := &model.InternalTransfer{
		ID:              uuid.New(),
		Coin:            "BTC",
		Amount:          decimal.NewFromFloat(1.0),
		FromAccountID:   fromAccountID,
		ToAccountID:     2,
		FromAccountType: "spot",
		ToAccountType:   "fund",
		Status:          model.ItsPending,
		UpdatedAt:       time.Now().UTC(),
	}
	transfer2 := &model.InternalTransfer{
		ID:              uuid.New(),
		Coin:            "ETH",
		Amount:          decimal.NewFromFloat(2.0),
		FromAccountID:   fromAccountID,
		ToAccountID:     545454,
		FromAccountType: "spot12",
		ToAccountType:   "fund12",
		Status:          model.ItsFailed,
		UpdatedAt:       time.Now().UTC().Add(time.Hour),
	}

	err := transferRepo.Add(ctx, transfer)
	require.NoError(t, err)
	err = transferRepo.Add(ctx, transfer2)
	require.NoError(t, err)

	retrievedTransfers, err := transferRepo.Find(ctx, &model.InternalTransferFilter{
		FromAccountID: &fromAccountID,
	})
	require.NoError(t, err)

	var count int
	for _, trans := range retrievedTransfers {
		switch trans.ID {
		case transfer.ID:
			require.Equal(t, transfer.ID, trans.ID)
			require.Equal(t, transfer.Coin, trans.Coin)
			require.True(t, transfer.Amount.Equal(trans.Amount))
			require.Equal(t, transfer.FromAccountID, trans.FromAccountID)
			require.Equal(t, transfer.ToAccountID, trans.ToAccountID)
			require.Equal(t, transfer.FromAccountType, trans.FromAccountType)
			require.Equal(t, transfer.ToAccountType, trans.ToAccountType)
			require.Equal(t, transfer.Status, trans.Status)
			require.WithinDuration(t, transfer.UpdatedAt, trans.UpdatedAt, time.Second)
			count++
		case transfer2.ID:
			require.Equal(t, transfer2.ID, trans.ID)
			require.Equal(t, transfer2.Coin, trans.Coin)
			require.True(t, transfer2.Amount.Equal(trans.Amount))
			require.Equal(t, transfer2.FromAccountID, trans.FromAccountID)
			require.Equal(t, transfer2.ToAccountID, trans.ToAccountID)
			require.Equal(t, transfer2.FromAccountType, trans.FromAccountType)
			require.Equal(t, transfer2.ToAccountType, trans.ToAccountType)
			require.Equal(t, transfer2.Status, trans.Status)
			require.WithinDuration(t, transfer2.UpdatedAt, trans.UpdatedAt, time.Second)
			count++
		}
	}
	require.Equal(t, 2, count)
}

func TestTransfer_AddFindOne(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	transferRepo := NewTransfer(db)

	transferID := uuid.New()
	transfer := &model.InternalTransfer{
		ID:              transferID,
		Coin:            "BTC",
		Amount:          decimal.NewFromFloat(1.0),
		FromAccountID:   5,
		ToAccountID:     2,
		FromAccountType: "spot",
		ToAccountType:   "fund",
		Status:          model.ItsPending,
		UpdatedAt:       time.Now().UTC(),
	}

	err := transferRepo.Add(ctx, transfer)
	require.NoError(t, err)

	retrievedTransfer, err := transferRepo.FindOne(ctx, &model.InternalTransferFilter{
		ID: &transferID,
	})
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, retrievedTransfer.ID)
	require.Equal(t, transfer.Coin, retrievedTransfer.Coin)
	require.True(t, transfer.Amount.Equal(retrievedTransfer.Amount))
	require.Equal(t, transfer.FromAccountID, retrievedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, retrievedTransfer.ToAccountID)
	require.Equal(t, transfer.FromAccountType, retrievedTransfer.FromAccountType)
	require.Equal(t, transfer.ToAccountType, retrievedTransfer.ToAccountType)
	require.Equal(t, transfer.Status, retrievedTransfer.Status)
	require.WithinDuration(t, transfer.UpdatedAt, retrievedTransfer.UpdatedAt, time.Second)
}

func TestTransfer_FindLastInternalTransfer(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	transferRepo := NewTransfer(db)

	transferID := uuid.New()
	fromAccountID := int64(1)
	transfer := &model.InternalTransfer{
		ID:              transferID,
		Coin:            "BTC",
		Amount:          decimal.NewFromFloat(1.0),
		FromAccountID:   fromAccountID,
		ToAccountID:     2,
		FromAccountType: "spot",
		ToAccountType:   "fund",
		Status:          model.ItsPending,
		UpdatedAt:       time.Now().UTC(),
	}

	transfer2 := &model.InternalTransfer{
		ID:              uuid.New(),
		Coin:            "ETH",
		Amount:          decimal.NewFromFloat(123),
		FromAccountID:   fromAccountID,
		ToAccountID:     5,
		FromAccountType: "fund",
		ToAccountType:   "spot",
		Status:          model.ItsSuccess,
		UpdatedAt:       time.Now().UTC().Add(-24 * time.Hour),
	}

	err := transferRepo.Add(ctx, transfer)
	require.NoError(t, err)

	err = transferRepo.Add(ctx, transfer2)
	require.NoError(t, err)

	isLast := true
	lastTransfer, err := transferRepo.FindOne(ctx, &model.InternalTransferFilter{
		FromAccountID: &fromAccountID,
		IsLast:        &isLast,
	})
	require.NoError(t, err)

	require.Equal(t, transfer.ID, lastTransfer.ID)
}

func TestTransfer_UpdateInternalTransferStatus(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	transferRepo := NewTransfer(db)

	transferID := uuid.New()
	transfer := &model.InternalTransfer{
		ID:              transferID,
		Coin:            "BTC",
		Amount:          decimal.NewFromFloat(1.0),
		FromAccountID:   10,
		ToAccountID:     20,
		FromAccountType: "spot",
		ToAccountType:   "fund",
		Status:          model.ItsPending,
		UpdatedAt:       time.Now().UTC(),
	}

	err := transferRepo.Add(ctx, transfer)
	require.NoError(t, err)

	st := model.ItsSuccess
	err = transferRepo.Update(ctx, transfer,
		&model.InternalTransferFilter{
			ID: &transferID,
		}, &model.InternalTransferPartial{
			Status: &st,
		})
	require.NoError(t, err)

	var updatedStatus model.InternalTransferStatus
	err = db.QueryRow(ctx, `
		SELECT status
		FROM swap.internal_transfers
		WHERE id = $1`, transferID).Scan(&updatedStatus)
	require.NoError(t, err)
	require.Equal(t, model.ItsSuccess, updatedStatus)
}
