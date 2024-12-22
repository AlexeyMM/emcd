package repository_migration_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	accountingModel "code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository/repository_migration"
	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

func TestMigration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "emcd.users_accounts", "emcd.addresses", "emcd.autopay_addresses")

	repo := repository_migration.NewMigrationRepository(dbPool, dbPool)

	userId := int32(1)
	userUuid := uuid.New()

	insertUser(ctx, t, dbPool, userId, userUuid)

	coinId := int32(1)
	coinCode := "btc"

	lastAt1 := time.Now().Add(-100 * time.Hour).UTC().Truncate(time.Millisecond)
	lastAt2 := time.Now().Add(-50 * time.Hour).UTC().Truncate(time.Millisecond)

	addressStr1 := "address1"
	addressStr2 := "address2"
	autopaAddresStr := "autopay_address"

	userAccount1 := &accountingModel.UserAccount{
		ID:            0,
		UserID:        userId,
		CoinID:        coinId,
		AccountTypeID: enum.NewAccountTypeIdWrapper(enum.WalletAccountTypeID),
		Minpay:        0,
		Address:       sql.NullString{String: addressStr1, Valid: true},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      sql.NullBool{Bool: true, Valid: true},
		CreatedAt:     sql.NullTime{Time: lastAt1.Add(time.Hour), Valid: true},
		UpdatedAt:     sql.NullTime{},
		Fee:           sql.NullFloat64{},
		UserIDNew:     uuid.NullUUID{},
		CoinNew:       sql.NullString{},
	}

	userAccount2 := &accountingModel.UserAccount{
		ID:            1,
		UserID:        userId,
		CoinID:        coinId,
		AccountTypeID: enum.NewAccountTypeIdWrapper(enum.WalletAccountTypeID),
		Minpay:        0,
		Address:       sql.NullString{String: addressStr2, Valid: true},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      sql.NullBool{Bool: true, Valid: true},
		CreatedAt:     sql.NullTime{Time: lastAt2.Add(time.Hour), Valid: true},
		UpdatedAt:     sql.NullTime{},
		Fee:           sql.NullFloat64{},
		UserIDNew:     uuid.NullUUID{UUID: userUuid, Valid: true},
		CoinNew:       sql.NullString{String: coinCode, Valid: true},
	}

	insertUsersAccount(ctx, t, dbPool, userAccount1)
	insertUsersAccount(ctx, t, dbPool, userAccount2)

	addressOffset := int32(1)

	address1 := &model_migration.AddressMigration{
		Id:            0,
		UserAccountId: userAccount1.ID,
		CoinId:        coinId,
		TokenId:       sql.NullInt32{},
		Address:       addressStr1,
		NetworkId:     sql.NullString{},
		CreatedAt:     lastAt1.Add(time.Hour),
		DeletedAt:     sql.NullTime{},
		AddressOffset: sql.NullInt32{},
		UserUuid:      uuid.UUID{},
	}

	address2 := &model_migration.AddressMigration{
		Id:            1,
		UserAccountId: userAccount2.ID,
		CoinId:        coinId,
		TokenId:       sql.NullInt32{},
		Address:       addressStr2,
		NetworkId:     sql.NullString{},
		CreatedAt:     lastAt2.Add(time.Hour),
		DeletedAt:     sql.NullTime{},
		AddressOffset: sql.NullInt32{Int32: addressOffset, Valid: true},
		UserUuid:      uuid.UUID{},
	}

	insertAddress(ctx, t, dbPool, address1)
	insertAddress(ctx, t, dbPool, address2)

	autopayAddress := &model_migration.AutopayAddressMigration{
		Id:            0,
		UserAccountId: 1,
		Address:       autopaAddresStr,
		Percent:       0,
		Label:         sql.NullString{},
		CreatedAt:     lastAt1,
		UpdatedAt:     lastAt1,
	}

	insertAutopayAddress(ctx, t, dbPool, autopayAddress)

	t.Run("get user accounts 1 migration success", func(t *testing.T) {
		totalCount2 := uint64(2)

		totalCount, userAccounts, err := repo.GetUserAccountMigrations(ctx, lastAt1, 1)
		require.NotNil(t, totalCount)
		require.NotEmpty(t, userAccounts)
		require.NoError(t, err)

		require.Len(t, userAccounts, 1)
		require.Equal(t, *totalCount, totalCount2)

		require.Equal(t, userAccounts[0].UserUuid, userUuid)

	})

	t.Run("get user accounts 2 migration success", func(t *testing.T) {
		totalCount1 := uint64(1)

		totalCount, userAccounts, err := repo.GetUserAccountMigrations(ctx, lastAt2, 1)
		require.NotNil(t, totalCount)
		require.NotEmpty(t, userAccounts)
		require.NoError(t, err)

		require.Len(t, userAccounts, 1)

		require.Equal(t, *totalCount, totalCount1)
		require.Equal(t, userAccounts[0].UserUuid, userUuid)
		require.Equal(t, userAccounts[0].Address.String, addressStr2)

	})

	t.Run("get address 1 migration success", func(t *testing.T) {
		totalCount2 := uint64(2)

		totalCount, addresses, err := repo.GetAddressMigrations(ctx, lastAt1, 1)
		require.NotNil(t, totalCount)
		require.NotEmpty(t, addresses)
		require.NoError(t, err)

		require.Len(t, addresses, 1)
		require.Equal(t, *totalCount, totalCount2)

		require.Equal(t, addresses[0].UserUuid, userUuid)

	})

	t.Run("get address 2 migration success", func(t *testing.T) {
		totalCount1 := uint64(1)

		totalCount, addresses, err := repo.GetAddressMigrations(ctx, lastAt2, 1)
		require.NotNil(t, totalCount)
		require.NotEmpty(t, addresses)
		require.NoError(t, err)

		require.Len(t, addresses, 1)

		require.Equal(t, *totalCount, totalCount1)
		require.Equal(t, addresses[0].UserUuid, userUuid)
		require.Equal(t, addresses[0].Address, addressStr2)

	})

	t.Run("get autopay address migration success", func(t *testing.T) {
		totalCount1 := uint64(1)

		totalCount, addresses, err := repo.GetAddressPersonalMigrations(ctx, lastAt1, 1)
		require.NotNil(t, totalCount)
		require.NotEmpty(t, addresses)
		require.NoError(t, err)

		require.Len(t, addresses, 1)

		require.Equal(t, *totalCount, totalCount1)
		require.Equal(t, addresses[0].UserUuid, userUuid)
		require.Equal(t, addresses[0].AaAddress, autopaAddresStr)

	})
}
