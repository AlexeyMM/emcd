package repository_test

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBalance_Change(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				// require.NoError(t, err) // TODO: check it
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		transactionID, err := repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)
		require.NotZero(t, transactionID)

		transactionFromDB := getTransaction(ctx, dbPool, t, transactionID)
		require.Equal(t, transaction.Type, transactionFromDB.Type)
		require.Equal(t, transaction.CreatedAt.Truncate(time.Second), transactionFromDB.CreatedAt.Truncate(time.Second))
		require.Equal(t, transaction.SenderAccountID, transactionFromDB.SenderAccountID)
		require.Equal(t, transaction.ReceiverAccountID, transactionFromDB.ReceiverAccountID)
		require.Equal(t, transaction.CoinID, transactionFromDB.CoinID)
		require.Equal(t, transaction.Amount, transactionFromDB.Amount)
		require.Equal(t, transaction.Comment, transactionFromDB.Comment)
		require.Equal(t, transaction.Hash, transactionFromDB.Hash)
		require.Equal(t, transaction.ReceiverAddress, transactionFromDB.ReceiverAddress)
		require.Equal(t, transaction.ActionID, transactionFromDB.ActionID)

		operations := getOperations(ctx, dbPool, t, transactionID)
		require.Len(t, operations, 2)
		for _, op := range operations {
			require.Equal(t, op.CoinID, coin.ID)
			require.Equal(t, op.Type, transaction.Type)
			require.Equal(t, int(op.TransactionID), transactionID)
			require.Equal(t, op.CreatedAt.Truncate(time.Second), transaction.CreatedAt.Truncate(time.Second))

			if op.AccountID == transaction.SenderAccountID {
				require.Equal(t, op.Amount, transaction.Amount.Neg())
				require.Equal(t, op.AccountID, transaction.SenderAccountID)
			} else {
				require.Equal(t, op.Amount, transaction.Amount)
				require.Equal(t, op.AccountID, transaction.ReceiverAccountID)
			}
		}
	})

	t.Run("Fail: not enough balance", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "user1", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)

		require.Error(t, err)
		require.ErrorContains(t, err, "balance is less than the transfer amount")
	})

	t.Run("Fail: actionId is empty", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "",
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)

		require.Error(t, err)
		require.ErrorContains(t, err, "action_id is empty")
	})
}

func TestBalance_ChangeWithBlock(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		blockedTill := time.Now().Add(time.Hour).UTC()
		block := model.Block{
			UnblockToAccountID: userReceiver.ID,
			BlockedTill:        blockedTill,
		}

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = blockedTill

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		transactionID, err := repository.ChangeWithBlock(ctx, sqlTx, &transaction)

		require.NoError(t, err)
		require.NotZero(t, transactionID)
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		transactionFromDB := getTransaction(ctx, dbPool, t, transactionID)
		require.Equal(t, transaction.Type, transactionFromDB.Type)
		require.Equal(t, transaction.CreatedAt.Truncate(time.Second), transactionFromDB.CreatedAt.Truncate(time.Second))
		require.Equal(t, transaction.SenderAccountID, transactionFromDB.SenderAccountID)
		require.Equal(t, transaction.ReceiverAccountID, transactionFromDB.ReceiverAccountID)
		require.Equal(t, transaction.CoinID, transactionFromDB.CoinID)
		require.Equal(t, transaction.Amount, transactionFromDB.Amount)
		require.Equal(t, transaction.Comment, transactionFromDB.Comment)
		require.Equal(t, transaction.Hash, transactionFromDB.Hash)
		require.Equal(t, transaction.ReceiverAddress, transactionFromDB.ReceiverAddress)
		require.Equal(t, transaction.ActionID, transactionFromDB.ActionID)

		operations := getOperations(ctx, dbPool, t, transactionID)
		require.Len(t, operations, 2)
		for _, op := range operations {
			require.Equal(t, op.CoinID, coin.ID)
			require.Equal(t, op.Type, transaction.Type)
			require.Equal(t, int(op.TransactionID), transactionID)
			require.Equal(t, op.CreatedAt.Truncate(time.Second), transaction.CreatedAt.Truncate(time.Second))

			if op.AccountID == transaction.SenderAccountID {
				require.Equal(t, op.Amount, transaction.Amount.Neg())
				require.Equal(t, op.AccountID, transaction.SenderAccountID)
			} else {
				require.Equal(t, op.Amount, transaction.Amount)
				require.Equal(t, op.AccountID, transaction.ReceiverAccountID)
			}
		}

		blockFromDB := getBlock(ctx, dbPool, t)
		require.NotZero(t, blockFromDB.ID)
		require.Equal(t, int(blockFromDB.BlockTransactionID), transactionID)
		require.Equal(t, block.UnblockToAccountID, blockFromDB.UnblockToAccountID)
		require.Equal(t, block.BlockedTill.Truncate(time.Second), blockFromDB.BlockedTill.Truncate(time.Second))
	})

	t.Run("Fail: blocking time has passed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		blockedTill := time.Now().UTC().Add(time.Duration(-1) * time.Hour)
		block := model.Block{
			UnblockToAccountID: userReceiver.ID,
			BlockedTill:        blockedTill,
		}
		require.Zero(t, block.ID)
		require.Zero(t, block.BlockTransactionID)

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = blockedTill

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.ChangeWithBlock(ctx, sqlTx, &transaction)

		require.Error(t, err)
		require.ErrorContains(t, err, "blocking time has passed")
	})

	t.Run("Fail: not enough balance", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "user1", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "user2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		blockedTill := time.Now().UTC().Add(time.Hour)
		block := model.Block{
			UnblockToAccountID: userReceiver.ID,
			BlockedTill:        blockedTill,
		}
		require.Zero(t, block.ID)
		require.Zero(t, block.BlockTransactionID)

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = blockedTill

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.ChangeWithBlock(ctx, sqlTx, &transaction)

		require.Error(t, err)
		require.ErrorContains(t, err, "balance is less than the transfer amount")
	})
}

func TestBalance_ChangeWithUnblock(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		blockTransaction := model.Transaction{
			Type:              34,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              blockTransaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            blockTransaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   blockTransaction.SenderAccountID,
			ReceiverAccountID: blockTransaction.ReceiverAccountID,
			CoinID:            blockTransaction.CoinID,
			Type:              int(blockTransaction.Type),
			ReceiverAddress:   blockTransaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           blockTransaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		blockedTill := time.Now().Add(time.Hour).UTC()
		block := model.Block{
			UnblockToAccountID: userReceiver.ID,
			BlockedTill:        blockedTill,
		}

		blockTransaction.UnblockAccountId = userReceiver.ID
		blockTransaction.BlockedTill = blockedTill

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		blockTransactionID, err := repository.ChangeWithBlock(ctx, sqlTx, &blockTransaction)
		require.NoError(t, err)

		unblockTransaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userReceiver.ID,
			ReceiverAccountID: userSender.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		require.Equal(t, blockTransaction.ActionID, unblockTransaction.ActionID)

		unblockTransactionID, err := repository.ChangeWithUnblock(ctx, sqlTx, &unblockTransaction)
		require.NoError(t, err)
		require.NotZero(t, unblockTransactionID)
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		blockFromDB := getBlock(ctx, dbPool, t)
		require.NotZero(t, blockFromDB.ID)
		require.Equal(t, blockTransactionID, int(blockFromDB.BlockTransactionID))
		require.Equal(t, unblockTransactionID, int(blockFromDB.UnblockTransactionID))
		require.Equal(t, block.UnblockToAccountID, blockFromDB.UnblockToAccountID)
		require.Equal(t, block.BlockedTill.Truncate(time.Second), blockFromDB.BlockedTill.Truncate(time.Second))

		unblockTransactionFromDB := getTransaction(ctx, dbPool, t, unblockTransactionID)
		require.Equal(t, unblockTransaction.Type, unblockTransactionFromDB.Type)
		require.Equal(t, unblockTransaction.CreatedAt.Truncate(time.Second), unblockTransactionFromDB.CreatedAt.Truncate(time.Second))
		require.Equal(t, unblockTransaction.SenderAccountID, unblockTransactionFromDB.SenderAccountID)
		require.Equal(t, unblockTransaction.ReceiverAccountID, unblockTransactionFromDB.ReceiverAccountID)
		require.Equal(t, unblockTransaction.CoinID, unblockTransactionFromDB.CoinID)
		require.Equal(t, unblockTransaction.Amount, unblockTransactionFromDB.Amount)
		require.Equal(t, unblockTransaction.Comment, unblockTransactionFromDB.Comment)
		require.Equal(t, unblockTransaction.Hash, unblockTransactionFromDB.Hash)
		require.Equal(t, unblockTransaction.ReceiverAddress, unblockTransactionFromDB.ReceiverAddress)
		require.Equal(t, unblockTransaction.ActionID, unblockTransactionFromDB.ActionID)
		require.Equal(t, blockTransaction.ActionID, unblockTransactionFromDB.ActionID)

		operations := getOperations(ctx, dbPool, t, unblockTransactionID)
		require.Len(t, operations, 2)
		for _, op := range operations {
			require.Equal(t, op.CoinID, coin.ID)
			require.Equal(t, op.Type, unblockTransaction.Type)
			require.Equal(t, int(op.TransactionID), unblockTransactionID)
			require.Equal(t, op.CreatedAt.Truncate(time.Second), unblockTransaction.CreatedAt.Truncate(time.Second))

			if op.AccountID == unblockTransaction.SenderAccountID {
				require.Equal(t, op.Amount, unblockTransaction.Amount.Neg())
				require.Equal(t, op.AccountID, unblockTransaction.SenderAccountID)
			} else {
				require.Equal(t, op.Amount, unblockTransaction.Amount)
				require.Equal(t, op.AccountID, unblockTransaction.ReceiverAccountID)
			}
		}
	})
}

func TestBalance_View(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types",
			"emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		count := 10
		amount := decimal.NewFromFloat(1)
		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction := model.Transaction{
				ID:                0,
				Type:              1,
				CreatedAt:         time.Now().UTC(),
				SenderAccountID:   userSender.ID,
				ReceiverAccountID: userReceiver.ID,
				CoinID:            coin.ID,
				Amount:            amount,
				Comment:           "Hello",
				Hash:              "Hello",
				ReceiverAddress:   "Hello",
				ActionID:          actionID.String(),
			}

			action := Action{
				UUID:              transaction.ActionID,
				Status:            "1",
				CreatedAt:         time.Now().UTC(),
				Amount:            transaction.Amount,
				Fee:               decimal.NewFromInt(0),
				SenderAccountID:   transaction.SenderAccountID,
				ReceiverAccountID: transaction.ReceiverAccountID,
				CoinID:            transaction.CoinID,
				Type:              int(transaction.Type),
				ReceiverAddress:   transaction.ReceiverAddress,
				Comment:           "test action",
				TokenID:           transaction.TokenID,
				UserID:            userSender.ID,
				SendMessage:       false,
			}
			writeAction(ctx, dbPool, t, action)

			_, err = repository.Change(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		resultSender, err := repository.View(ctx, userSender.ID, coin.ID, uaSender.AccountTypeID.AccountTypeId, false)
		require.NoError(t, err)
		sum := amount.Neg().IntPart() * int64(count)
		require.Equal(t, sum, resultSender.IntPart())

		resultReceiver, err := repository.View(ctx, userReceiver.ID, coin.ID, uaReceiver.AccountTypeID.AccountTypeId, false)
		require.NoError(t, err)
		sum = amount.IntPart() * int64(count)
		require.Equal(t, sum, resultReceiver.IntPart())
	})

	t.Run("Success: with total balance", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types",
			"emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		var (
			count              = 2
			blockTransactionID int
			amount             = decimal.NewFromFloat(10)
		)
		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)
		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction := model.Transaction{
				Type:              1,
				CreatedAt:         time.Now().UTC(),
				SenderAccountID:   userSender.ID,
				ReceiverAccountID: userReceiver.ID,
				CoinID:            coin.ID,
				Amount:            amount,
				Comment:           "Hello",
				Hash:              "Hello",
				ReceiverAddress:   "Hello",
				ActionID:          actionID.String(),
			}

			action := Action{
				UUID:              transaction.ActionID,
				Status:            "1",
				CreatedAt:         time.Now().UTC(),
				Amount:            transaction.Amount,
				Fee:               decimal.NewFromInt(0),
				SenderAccountID:   transaction.SenderAccountID,
				ReceiverAccountID: transaction.ReceiverAccountID,
				CoinID:            transaction.CoinID,
				Type:              int(transaction.Type),
				ReceiverAddress:   transaction.ReceiverAddress,
				Comment:           "test action",
				TokenID:           transaction.TokenID,
				UserID:            userSender.ID,
				SendMessage:       false,
			}
			writeAction(ctx, dbPool, t, action)

			blockTransactionID, err = repository.Change(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		// add last transaction to transactions_blocks
		query := `
		INSERT INTO emcd.transactions_blocks
		    (block_transaction_id, unblock_to_account_id, blocked_till)
		VALUES ($1, $2, $3)`
		_, err = dbPool.Exec(ctx, query, blockTransactionID, 1, time.Now().Add(time.Minute))
		require.NoError(t, err)

		// check without totalBalance
		result, err := repository.View(ctx, userReceiver.ID, coin.ID, uaReceiver.AccountTypeID.AccountTypeId, false)
		require.NoError(t, err)
		require.Equal(t, int64(count)*amount.IntPart(), result.IntPart())

		// check with totalBalance
		result, err = repository.View(ctx, userReceiver.ID, coin.ID, uaReceiver.AccountTypeID.AccountTypeId, true)
		require.NoError(t, err)
		require.Equal(t, int64(count+1)*amount.IntPart(), result.IntPart())
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		_, err := repository.View(ctx, 1, 1, 1, false)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("Fail: not found with total balance", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		_, err := repository.View(ctx, 1, 1, 1, true)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestBalance_FindOperations(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		count := 10
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction.ActionID = actionID.String()

			action := Action{
				UUID:              transaction.ActionID,
				Status:            "1",
				CreatedAt:         time.Now().UTC(),
				Amount:            transaction.Amount,
				Fee:               decimal.NewFromInt(0),
				SenderAccountID:   transaction.SenderAccountID,
				ReceiverAccountID: transaction.ReceiverAccountID,
				CoinID:            transaction.CoinID,
				Type:              int(transaction.Type),
				ReceiverAddress:   transaction.ReceiverAddress,
				Comment:           "test action",
				TokenID:           transaction.TokenID,
				UserID:            userSender.ID,
				SendMessage:       false,
			}
			writeAction(ctx, dbPool, t, action)

			_, err = repository.Change(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		operations, err := repository.FindOperations(ctx, int(userReceiver.ID), int(coin.ID))
		require.NoError(t, err)
		require.Len(t, operations, count)
		for _, o := range operations {
			require.Equal(t, int64(uaReceiver.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaReceiver.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
			require.Equal(t, transaction.Amount, o.Amount)
			require.Equal(t, uaReceiver.IsActive.Bool, o.IsActive) // TODO: better both type is sql.NullBool
			require.NotZero(t, o.TransactionID)
		}

		operations, err = repository.FindOperations(ctx, int(userSender.ID), int(coin.ID))
		require.NoError(t, err)
		require.Len(t, operations, count)
		for _, o := range operations {
			require.Equal(t, int64(uaSender.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaSender.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
			require.Equal(t, transaction.Amount.Neg(), o.Amount)
			require.Equal(t, uaSender.IsActive.Bool, o.IsActive) // TODO: better both type is sql.NullBool
			require.NotZero(t, o.TransactionID)
		}
	})

	t.Run("Success: find with blocks", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
		}

		blockedTill := time.Now().Add(time.Hour).UTC()

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = blockedTill
		// require.Zero(t, block.ID)
		// require.Zero(t, block.BlockTransactionID)

		action := Action{
			Status:            "1",
			CreatedAt:         time.Now().UTC(),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		count := 10
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction.ActionID = actionID.String()

			action.UUID = actionID.String()
			writeAction(ctx, dbPool, t, action)

			_, err = repository.ChangeWithBlock(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		operations, err := repository.FindOperations(ctx, int(userReceiver.ID), int(coin.ID))
		require.NoError(t, err)
		require.Len(t, operations, count)
		for _, o := range operations {
			require.Equal(t, int64(uaReceiver.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaReceiver.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
			require.Equal(t, transaction.Amount, o.Amount)
			require.Equal(t, uaReceiver.IsActive.Bool, o.IsActive) // TODO: better both type is sql.NullBool
			require.NotZero(t, o.UnblockToAccountID)
			require.NotZero(t, o.TransactionID)
		}

		operations, err = repository.FindOperations(ctx, int(userSender.ID), int(coin.ID))
		require.NoError(t, err)
		require.Len(t, operations, count)
		for _, o := range operations {
			require.Equal(t, int64(uaSender.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaSender.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
			require.Equal(t, transaction.Amount.Neg(), o.Amount)
			require.Equal(t, uaSender.IsActive.Bool, o.IsActive) // TODO: better both type is sql.NullBool
			require.NotZero(t, o.UnblockToAccountID)
			require.NotZero(t, o.TransactionID)
		}
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		operations, err := repository.FindOperations(ctx, 1, 1)
		require.NoError(t, err)
		require.Empty(t, operations)
	})
}

func TestBalance_FindBatchOperations(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			ID:                0,
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
		}
		require.Zero(t, transaction.ID)

		action := Action{
			Status:            "1",
			CreatedAt:         time.Now().UTC(),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		count := 10
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction.ActionID = actionID.String()

			action.UUID = actionID.String()
			writeAction(ctx, dbPool, t, action)

			_, err = repository.Change(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		result, err := repository.FindBatchOperations(
			ctx,
			map[int]int{
				int(userSender.ID):   int(coin.ID),
				int(userReceiver.ID): int(coin.ID),
			},
		)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Len(t, result[int(uaReceiver.ID)], count)
		require.Len(t, result[int(uaSender.ID)], count)

		for _, o := range result[int(uaReceiver.ID)] {
			require.Equal(t, int64(uaReceiver.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaReceiver.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
		}
		for _, o := range result[int(uaSender.ID)] {
			require.Equal(t, int64(uaSender.ID), o.AccountID)
			require.Equal(t, transaction.CoinID, o.OperationCoinID)
			require.Equal(t, transaction.CoinID, o.UserAccountCoinID)
			require.Equal(t, uaSender.AccountTypeID.AccountTypeId, o.AccountTypeID)
			require.Equal(t, transaction.Type, o.Type)
			require.Equal(t, transaction.CreatedAt.Truncate(time.Second), o.CreatedAt.Truncate(time.Second))
		}
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		operations, err := repository.FindBatchOperations(ctx, map[int]int{1: 1, 2: 2})
		require.NoError(t, err)
		require.Empty(t, operations)
	})
}

func TestBalance_FindTransactions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			ID:                0,
			Type:              model.MainCoinMiningPayoutTrTypeID,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "my-action-id",
		}

		action := Action{
			Status:            "1",
			CreatedAt:         time.Now().UTC(),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		count := 10
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction.ActionID = actionID.String()

			action.UUID = transaction.ActionID
			writeAction(ctx, dbPool, t, action)

			_, err = repository.Change(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		transactions, err := repository.FindTransactions(
			ctx,
			[]int{int(model.MainCoinMiningPayoutTrTypeID), int(model.MergeCoinMiningPayoutTrTypeID)},
			int(userReceiver.ID),
			accountType.ToInt(),
			[]int{int(coin.ID), 777},
			time.Now().Add(-(time.Hour * 24)),
		)
		require.NoError(t, err)
		require.Len(t, transactions, count)

	})

	t.Run("Success: passing zero time", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              model.MainCoinMiningPayoutTrTypeID,
			CreatedAt:         time.Now().UTC().Truncate(time.Millisecond),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          uuid.New().String(),
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		transactions, err := repository.FindTransactions(
			ctx,
			[]int{int(model.MainCoinMiningPayoutTrTypeID), int(model.MergeCoinMiningPayoutTrTypeID)},
			int(userReceiver.ID),
			accountType.ToInt(),
			[]int{int(coin.ID), 777},
			time.Time{},
		)
		require.NoError(t, err)
		require.Len(t, transactions, 1)

		result := transactions[0]
		require.Equal(t, transaction.Type, result.Type)
		require.Equal(t, transaction.CreatedAt.String(), result.CreatedAt.String())
		require.Equal(t, transaction.SenderAccountID, result.SenderAccountID)
		require.Equal(t, transaction.ReceiverAccountID, result.ReceiverAccountID)
		require.Equal(t, transaction.CoinID, result.CoinID)
		require.Equal(t, transaction.Amount, result.Amount)
		require.Equal(t, transaction.Comment, result.Comment)
		require.Equal(t, transaction.Hash, result.Hash)
		require.Equal(t, transaction.ReceiverAddress, result.ReceiverAddress)
		require.Equal(t, transaction.ActionID, result.ActionID)
	})

	t.Run("Success: null values", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              model.MainCoinMiningPayoutTrTypeID,
			CreatedAt:         time.Now().UTC().Truncate(time.Millisecond),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			ActionID:          uuid.New().String(),
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		transactions, err := repository.FindTransactions(
			ctx,
			[]int{int(model.MainCoinMiningPayoutTrTypeID), int(model.MergeCoinMiningPayoutTrTypeID)},
			int(userReceiver.ID),
			accountType.ToInt(),
			[]int{int(coin.ID), 777},
			time.Time{},
		)
		require.NoError(t, err)
		require.Len(t, transactions, 1)

		result := transactions[0]
		require.Equal(t, transaction.Type, result.Type)
		require.Equal(t, transaction.CreatedAt.String(), result.CreatedAt.String())
		require.Equal(t, transaction.SenderAccountID, result.SenderAccountID)
		require.Equal(t, transaction.ReceiverAccountID, result.ReceiverAccountID)
		require.Equal(t, transaction.CoinID, result.CoinID)
		require.Equal(t, transaction.Amount, result.Amount)
		require.Empty(t, result.Comment)
		require.Empty(t, result.Hash)
		require.Empty(t, result.ReceiverAddress)
		require.Equal(t, transaction.ActionID, result.ActionID)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		transactions, err := repository.FindTransactions(ctx, []int{1}, 1, 1, []int{1}, time.Now())
		require.NoError(t, err)
		require.Empty(t, transactions)
	})
}

func TestBalance_FindTransactionsWithBlocks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			ID:                0,
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
		}

		block := model.Block{
			ID:                   0,
			BlockTransactionID:   0,
			UnblockTransactionID: 0,
			UnblockToAccountID:   userReceiver.ID,
			BlockedTill:          time.Now().Add(time.Hour),
		}
		require.Zero(t, block.ID)
		require.Zero(t, block.BlockTransactionID)
		blockedTill := time.Now().Add(time.Hour * 2)
		require.True(t, block.BlockedTill.Before(blockedTill))

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = time.Now().Add(time.Hour)

		action := Action{
			Status:            "1",
			CreatedAt:         time.Now().UTC(),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		count := 10
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		for i := 0; i < count; i++ {
			actionID := uuid.New()
			transaction.ActionID = actionID.String()

			action.UUID = transaction.ActionID
			writeAction(ctx, dbPool, t, action)

			_, err = repository.ChangeWithBlock(ctx, sqlTx, &transaction)
			require.NoError(t, err)
		}
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		transactions, err := repository.FindTransactionsWithBlocks(ctx, blockedTill)
		require.NoError(t, err)
		require.Len(t, transactions, count)
		for _, v := range transactions {
			require.Equal(t, transaction.ReceiverAccountID, v.ReceiverAccountID)
			require.Equal(t, transaction.CoinID, v.CoinID)
			require.Equal(t, transaction.Type, v.Type)
			require.Equal(t, transaction.Amount, v.Amount)
			require.NotZero(t, v.BlockID)
			require.Equal(t, block.UnblockToAccountID, v.UnblockToAccountID)
		}
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		transactions, err := repository.FindTransactionsWithBlocks(ctx, time.Now())
		require.NoError(t, err)
		require.Empty(t, transactions)
	})
}

func TestBalance_GetTransactionByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.transactions_blocks", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			ID:                0,
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		block := model.Block{
			UnblockToAccountID: userReceiver.ID,
			BlockedTill:        time.Now().Add(time.Hour),
		}

		transaction.UnblockAccountId = userReceiver.ID
		transaction.BlockedTill = time.Now().Add(time.Hour)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		transactionID, err := repository.ChangeWithBlock(ctx, sqlTx, &transaction)
		require.NoError(t, err)
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		result, err := repository.GetTransactionByID(ctx, transactionID)
		require.NoError(t, err)
		require.Equal(t, transaction.SenderAccountID, result.SenderAccountID)
		require.Equal(t, transaction.ReceiverAccountID, result.ReceiverAccountID)
		require.Equal(t, block.UnblockToAccountID, result.UnblockToAccountID)
		require.Equal(t, transaction.CoinID, result.CoinID)
		require.Equal(t, transaction.Amount, result.Amount)
		require.Equal(t, transaction.Type, result.Type)
		require.Zero(t, result.UnblockTransactionID)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		transaction, err := repository.GetTransactionByID(ctx, 1)
		require.ErrorIs(t, err, pgx.ErrNoRows)
		require.Nil(t, transaction)
	})
}

func TestBalance_FindLastBlockTimeBalances(t *testing.T) {
	t.Run("Success: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		data := []int64{1}
		result, err := repository.FindLastBlockTimeBalances(ctx, data)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)
		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			ID:                0,
			Type:              model.IncomeBillTrTypeID,
			CreatedAt:         time.Now().Add(-time.Hour),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		transaction = model.Transaction{
			ID:                0,
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         time.Now(),
			SenderAccountID:   userReceiver.ID,
			ReceiverAccountID: userSender.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		data := []int64{int64(uaSender.ID), int64(uaReceiver.ID)}
		result, err := repository.FindLastBlockTimeBalances(ctx, data)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.True(t, result[int(uaReceiver.ID)].IsZero())
		require.True(t, result[int(uaSender.ID)].IsZero())
	})
}

func TestBalance_FindBalancesDiffMining(t *testing.T) {
	t.Run("Success: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		data := []*model.UserBeforePayoutMining{{UserID: 1, CoinID: 1, BlockID: 1, AccountTypeID: 1, LastPay: time.Now()}}
		result, err := repository.FindBalancesDiffMining(ctx, data)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountTypeMining := enum.AccountTypeIdWrapper{AccountTypeId: enum.MiningAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeMining)
		accountTypeWallet := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeWallet)
		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeWallet, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeMining, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		lastPay := time.Now().Add(-time.Hour * 10)

		transaction := model.Transaction{
			ID:                0,
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(10),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		blockTransaction := model.Transaction{
			ID:                0,
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour * 2),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		blockTransactionID, err := repository.Change(ctx, sqlTx, &blockTransaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		data := []*model.UserBeforePayoutMining{
			{UserID: userReceiver.ID, CoinID: coin.ID, BlockID: int64(blockTransactionID), AccountTypeID: accountTypeMining.AccountTypeId, LastPay: lastPay},
		}
		result, err := repository.FindBalancesDiffMining(ctx, data)
		require.NoError(t, err)
		require.Equal(t, result[blockTransactionID], transaction.Amount.Sub(blockTransaction.Amount))
	})

	t.Run("Success: no last pay", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountTypeMining := enum.AccountTypeIdWrapper{AccountTypeId: enum.MiningAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeMining)
		accountTypeWallet := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeWallet)
		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeWallet, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeMining, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		lastPay := time.Now().Add(-time.Hour * 10)

		transaction := model.Transaction{
			ID:                0,
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(10),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		blockTransaction := model.Transaction{
			ID:                0,
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour * 2),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		blockTransactionID, err := repository.Change(ctx, sqlTx, &blockTransaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		data := []*model.UserBeforePayoutMining{
			{UserID: userReceiver.ID, CoinID: coin.ID, BlockID: int64(blockTransactionID), AccountTypeID: accountTypeMining.AccountTypeId},
		}
		require.Zero(t, data[0].LastPay)

		result, err := repository.FindBalancesDiffMining(ctx, data)
		require.NoError(t, err)
		require.Equal(t, result[blockTransactionID], transaction.Amount.Sub(blockTransaction.Amount))
	})
}

func TestBalance_FindBalancesDiffWallet(t *testing.T) {
	t.Run("Success: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		data := []*model.UserBeforePayoutWallet{{UserID: 1, CoinID: 1, TransactionIDs: []int64{1}, AccountTypeID: 1}}
		result, err := repository.FindBalancesDiffWallet(ctx, data)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountTypeMining := enum.AccountTypeIdWrapper{AccountTypeId: enum.MiningAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeMining)
		accountTypeWallet := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeWallet)
		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "2", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeMining, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeWallet, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		lastPay := time.Now().Add(-time.Hour * 10)

		transaction := model.Transaction{
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(10),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		action := Action{
			UUID:              transaction.ActionID,
			Status:            "1",
			CreatedAt:         time.Now().UTC().Add(-1 * time.Minute),
			Amount:            transaction.Amount,
			Fee:               decimal.NewFromInt(0),
			SenderAccountID:   transaction.SenderAccountID,
			ReceiverAccountID: transaction.ReceiverAccountID,
			CoinID:            transaction.CoinID,
			Type:              int(transaction.Type),
			ReceiverAddress:   transaction.ReceiverAddress,
			Comment:           "test action",
			TokenID:           transaction.TokenID,
			UserID:            userSender.ID,
			SendMessage:       false,
		}
		writeAction(ctx, dbPool, t, action)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		_, err = repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)

		blockTransaction := model.Transaction{
			Type:              model.PoolPaysUsersBalanceTrTypeID,
			CreatedAt:         lastPay.Add(time.Hour * 2),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		blockTransactionID, err := repository.Change(ctx, sqlTx, &blockTransaction)
		require.NoError(t, err)

		err = sqlTx.Commit(ctx)
		require.NoError(t, err)

		data := []*model.UserBeforePayoutWallet{
			{UserID: userReceiver.ID, CoinID: coin.ID, TransactionIDs: []int64{int64(blockTransactionID)}, AccountTypeID: accountTypeWallet.AccountTypeId},
		}
		result, err := repository.FindBalancesDiffWallet(ctx, data)
		require.NoError(t, err)
		require.Equal(t, result[0].Diff, transaction.Amount.Sub(blockTransaction.Amount))
		require.Equal(t, result[0].BlockID, int64(blockTransactionID))
		require.Equal(t, result[0].UserID, userReceiver.ID)
	})
}

func TestBalance_GetAccountTypeIDByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountTypeMining := enum.AccountTypeIdWrapper{AccountTypeId: enum.MiningAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeMining)
		user := model.User{ID: 1, Username: "1", Password: "1"}
		writeUser(ctx, dbPool, t, user)
		ua := model.UserAccount{ID: 1, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeMining, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, ua)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				require.NoError(t, err)

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetAccountTypeIDByID(ctx, sqlTx, int(ua.ID))
		require.NoError(t, err)
		require.Equal(t, result, accountTypeMining.AccountTypeId)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetAccountTypeIDByID(ctx, sqlTx, 777)
		require.Zero(t, result)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestBalance_GetBlockAccountIDFor31Type(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)
		accountTypeBlock := enum.AccountTypeIdWrapper{AccountTypeId: enum.BlockUserAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeBlock)
		user := model.User{ID: 1, Username: "payouts_node", Password: "1"}
		writeUser(ctx, dbPool, t, user)
		ua := model.UserAccount{ID: 1, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeBlock, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, ua)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDFor31Type(ctx, sqlTx, int(coin.ID), user.Username)
		require.NoError(t, err)
		require.Equal(t, int32(result), ua.ID) // TODO: id user account is int32
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDFor31Type(ctx, sqlTx, 777, "Hello")
		require.Zero(t, result)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestBalance_GetBlockAccountIDFor57Type(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountTypeWallet := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeWallet)
		accountTypeBlock := enum.AccountTypeIdWrapper{AccountTypeId: enum.BlockUserAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeBlock)

		user := model.User{ID: 1, Username: "emcd_kucoin", Password: "1"}
		writeUser(ctx, dbPool, t, user)

		uaWallet := model.UserAccount{ID: 1, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeWallet, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaWallet)
		uaBlock := model.UserAccount{ID: 2, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeBlock, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaBlock)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDFor57Type(ctx, sqlTx, int(coin.ID), user.Username)
		require.NoError(t, err)
		require.Equal(t, result, uaBlock.ID)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDFor57Type(ctx, sqlTx, 777, "Hello")
		require.Zero(t, result)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestBalance_GetBlockAccountIDBySenderAccountID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations",
			"emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountTypeWallet := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeWallet)
		accountTypeBlock := enum.AccountTypeIdWrapper{AccountTypeId: enum.BlockUserAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountTypeBlock)

		user := model.User{ID: 1, Username: "emcd_kucoin", Password: "1"}
		writeUser(ctx, dbPool, t, user)

		uaWallet := model.UserAccount{ID: 1, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeWallet, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaWallet)
		uaBlock := model.UserAccount{ID: 2, UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountTypeBlock, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaBlock)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDBySenderAccountID(ctx, sqlTx, int(uaWallet.ID))
		require.NoError(t, err)
		require.Equal(t, result, uaBlock.ID)
	})

	t.Run("Fail: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		result, err := repository.GetBlockAccountIDBySenderAccountID(ctx, sqlTx, 777)
		require.Zero(t, result)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestBalance_GetTransactionIDByAction(t *testing.T) {
	t.Run("Success: not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)
		actionID := uuid.New()
		result, err := repository.GetTransactionIDByAction(ctx, actionID.String(), 1, "1")
		require.ErrorIs(t, err, pgx.ErrNoRows)
		require.Zero(t, result)
	})

	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC(),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            decimal.NewFromInt(1),
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}

		whiteList := []string{"admin"}
		repository := repository.NewBalance(dbPool, whiteList)

		sqlTx, err := dbPool.Begin(ctx)
		require.NoError(t, err)
		defer func(sqlTx pgx.Tx, ctx context.Context) {
			if err := sqlTx.Rollback(ctx); err != nil {
				sdkLog.Error(ctx, err.Error())

			}
		}(sqlTx, ctx)

		transactionID, err := repository.Change(ctx, sqlTx, &transaction)
		require.NoError(t, err)
		err = sqlTx.Commit(ctx)
		require.NoError(t, err)
		require.NotZero(t, transactionID)

		result, err := repository.GetTransactionIDByAction(ctx, transaction.ActionID, int(transaction.Type), transaction.Amount.String())
		require.NoError(t, err)
		require.Equal(t, int(result), transactionID)
	})
}

func TestBalance_GetUserIDsByAccountID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.users_accounts")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		oldID := int64(1)
		newID := uuid.New()
		refID := 444
		user := model.User{ID: oldID, Username: "admin", Password: "1", NewID: newID, RefID: refID}
		writeUserWithNewID(ctx, dbPool, t, user)
		accID := 2
		userAccount := model.UserAccount{ID: int32(accID), UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, userAccount)
		repository := repository.NewBalance(dbPool, nil)
		actualNew, actualOld, actualRefID, err := repository.GetUserIDsByAccountID(ctx, accID)
		require.NoError(t, err)
		require.Equal(t, newID, actualNew)
		require.Equal(t, int32(oldID), actualOld)
		require.Equal(t, int32(refID), actualRefID)
	})
}

func TestBalance_GetUserAccountIDByOldID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.users_accounts")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.MiningAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		oldID := int64(1)
		newID := uuid.New()
		user := model.User{ID: oldID, Username: "admin", Password: "1", NewID: newID}
		writeUserWithNewID(ctx, dbPool, t, user)
		accID := 2
		userAccount := model.UserAccount{ID: int32(accID), UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, userAccount)
		repository := repository.NewBalance(dbPool, nil)
		actualAccID, err := repository.GetUserAccountIDByOldID(ctx, int32(oldID), enum.MiningAccountTypeID, int(coin.ID))
		require.NoError(t, err)
		require.Equal(t, accID, actualAccID)
	})
}

func TestBalance_GetUserAccountIDByNewID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.users_accounts")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.ReferralAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		user := model.User{ID: 2, Username: "admin", Password: "1", NewID: uuid.New()}
		writeUserWithNewID(ctx, dbPool, t, user)
		accID := 2
		userAccount := model.UserAccount{ID: int32(accID), UserID: int32(user.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, userAccount)
		repository := repository.NewBalance(dbPool, nil)
		actualAccID, err := repository.GetUserAccountIDByNewID(ctx, user.NewID, enum.ReferralAccountTypeID, int(coin.ID))
		require.NoError(t, err)
		require.Equal(t, accID, actualAccID)
	})
}
