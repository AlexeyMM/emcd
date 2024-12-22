package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/accounting/model"
)

func TestReferralStatistic(t *testing.T) {
	t.Run("GetTotalLastMonthSimple", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome{{CoinID: coin.ID, Amount: amount.InexactFloat64()}}, resp.Income.Yesterday)
		require.Equal(t, []*model.ReferralIncome{{CoinID: coin.ID, Amount: amount.InexactFloat64()}}, resp.Income.ThisMonth)
	})
	t.Run("GetThisMonthOnly", func(t *testing.T) {
		if time.Now().Day() <= 2 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -2),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome{{CoinID: coin.ID, Amount: amount.InexactFloat64()}}, resp.Income.ThisMonth)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
	})
	t.Run("GetNothingNoTransactionsThatRange", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, -1, -1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.ThisMonth)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
	})
	t.Run("GetThisMonthSeveralCoins", func(t *testing.T) {
		if time.Now().Day() <= 3 { // TODO: make more case for tests

			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer truncateTables(
			ctx, t, "emcd.coins", "emcd.users", "emcd.account_types", "emcd.transactions", "emcd.operations", "emcd.users_accounts", "emcd.wallets_accounting_actions")

		coin := model.Coin{ID: 1, Name: "1", Description: "1", Code: "1", Rate: 1.1}
		writeCoin(ctx, dbPool, t, coin)

		coin2 := model.Coin{ID: 2, Name: "2", Description: "2", Code: "2", Rate: 2.2}
		writeCoin(ctx, dbPool, t, coin2)

		accountType := enum.AccountTypeIdWrapper{AccountTypeId: enum.WalletAccountTypeID}
		writeAccountType(ctx, dbPool, t, accountType)

		userSender := model.User{ID: 1, Username: "admin", Password: "1"}
		writeUser(ctx, dbPool, t, userSender)
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		uaSender2 := model.UserAccount{ID: 3, UserID: int32(userSender.ID), CoinID: int32(coin2.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender2)
		uaReceiver2 := model.UserAccount{ID: 4, UserID: int32(userReceiver.ID), CoinID: int32(coin2.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver2)

		amount := decimal.NewFromInt(1)
		transaction1 := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -2),
			SenderAccountID:   int64(uaSender.ID),
			ReceiverAccountID: int64(uaReceiver.ID),
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction1)
		transaction2 := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -3),
			SenderAccountID:   int64(uaSender2.ID),
			ReceiverAccountID: int64(uaReceiver2.ID),
			CoinID:            coin2.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction2)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(uaReceiver.ID), int64(uaReceiver2.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome{{CoinID: coin.ID, Amount: amount.InexactFloat64()}, {CoinID: coin2.ID, Amount: amount.InexactFloat64()}}, resp.Income.ThisMonth)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)

	})
	t.Run("GetNothingNoRequiredTransactions", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.ThisMonth)
	})
	t.Run("GetNothingNoTransactionsAtAll", func(t *testing.T) {
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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.ThisMonth)
	})
	t.Run("GetNothingNoTransactionForUserID", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)
		uaReceiver2 := model.UserAccount{ID: 100, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver2)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              22,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: int64(uaReceiver2.ID),
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              uuid.New(),
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.ThisMonth)
	})
	t.Run("GetErrorNoAccountIDs", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC().AddDate(0, -1, 1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              uuid.New(),
			TransactionTypesIDs: []int64{22},
		}
		_, err := repository.GetReferralsStatistic(ctx, &req)
		require.Error(t, err)
	})
	t.Run("GetNothingNoTransactionsForUserAccounts", func(t *testing.T) {
		if time.Now().Day() <= 1 { // TODO: make more case for tests

			return
		}

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
		userReceiver := model.User{ID: 2, NewID: uuid.New(), Username: "user", Password: "2"}
		writeUser(ctx, dbPool, t, userReceiver)

		uaSender := model.UserAccount{ID: 1, UserID: int32(userSender.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaSender)
		uaReceiver := model.UserAccount{ID: 2, UserID: int32(userReceiver.ID), CoinID: int32(coin.ID), AccountTypeID: accountType, Minpay: 1.1, IsActive: sql.NullBool{Bool: true, Valid: true}}
		writeUsersAccounts(ctx, dbPool, t, uaReceiver)

		amount := decimal.NewFromInt(1)
		transaction := model.Transaction{
			Type:              1,
			CreatedAt:         time.Now().UTC().AddDate(0, 0, -1),
			SenderAccountID:   userSender.ID,
			ReceiverAccountID: userReceiver.ID,
			CoinID:            coin.ID,
			Amount:            amount,
			Comment:           "Hello",
			Hash:              "Hello",
			ReceiverAddress:   "Hello",
			ActionID:          "8597ab20-257c-408e-be8e-f9b150bbe4b9",
		}
		writeTransaction(ctx, dbPool, t, &transaction)
		repository := repository.NewReferralStatistic(dbPool)
		req := model.ReferralsStatisticInput{
			UserID:              userReceiver.NewID,
			AccountIDs:          []int64{int64(userReceiver.ID)},
			TransactionTypesIDs: []int64{22},
		}
		resp, err := repository.GetReferralsStatistic(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.Yesterday)
		require.Equal(t, []*model.ReferralIncome(nil), resp.Income.ThisMonth)
	})
}

func writeTransaction(ctx context.Context, pool *pgxpool.Pool, t *testing.T, trx *model.Transaction) {
	query := `insert into emcd.transactions (type, sender_account_id, receiver_account_id, coin_id, amount, created_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := pool.Exec(ctx, query, trx.Type, trx.SenderAccountID, trx.ReceiverAccountID, trx.CoinID, trx.Amount, trx.CreatedAt)
	require.NoError(t, err)
}
