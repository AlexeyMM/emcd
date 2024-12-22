package repository_migration_test

import (
	"context"
	"fmt"
	"testing"

	accountingModel "code.emcdtech.com/emcd/service/accounting/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

func truncateTables(ctx context.Context, t *testing.T, dbPool *pgxpool.Pool, tables ...string) {
	for _, table := range tables {
		_, err := dbPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s cascade", table))
		require.NoError(t, err)
	}
}

func insertUser(ctx context.Context, t *testing.T, pool *pgxpool.Pool, userId int32, userUuid uuid.UUID) {
	query := `insert into emcd.users (id, new_id, username, password) values ($1,$2,$3,$4);`
	userName, userPass := "", ""
	_, err := pool.Exec(ctx, query, userId, userUuid, userName, userPass)

	require.NoError(t, err)
}

func insertUsersAccount(ctx context.Context, t *testing.T, pool *pgxpool.Pool, ua *accountingModel.UserAccount) {
	query := `insert into emcd.users_accounts (id, user_id, coin_id, account_type_id, address, minpay, coin_new, user_id_new, is_active, created_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`
	_, err := pool.Exec(ctx, query, ua.ID, ua.UserID, ua.CoinID, ua.AccountTypeID.ToInt32(), ua.Address, ua.Minpay, ua.CoinNew, ua.UserIDNew, ua.IsActive, ua.CreatedAt)

	require.NoError(t, err)
}

func insertAddress(ctx context.Context, t *testing.T, pool *pgxpool.Pool, a *model_migration.AddressMigration) {
	query := `insert into emcd.addresses (id, user_account_id, coin_id, token_id, address, created_at, address_offset) values ($1,$2,$3,$4,$5,$6,$7);`
	_, err := pool.Exec(ctx, query, a.Id, a.UserAccountId, a.CoinId, a.TokenId, a.Address, a.CreatedAt, a.AddressOffset)

	require.NoError(t, err)
}

func insertAutopayAddress(ctx context.Context, t *testing.T, pool *pgxpool.Pool, a *model_migration.AutopayAddressMigration) {
	query := `insert into emcd.autopay_addresses (id, user_account_id, address, created_at) values ($1,$2,$3,$4);`
	_, err := pool.Exec(ctx, query, a.Id, a.UserAccountId, a.Address, a.CreatedAt)

	require.NoError(t, err)
}
