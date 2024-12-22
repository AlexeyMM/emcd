package repository_migration

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

const sqlQuery = `SELECT
	ua.user_id_new as user_uuid,
	COALESCE(ua.coin_new, '') as coin_code,
	ua.minpay as min_payout,
	COALESCE(ua.address, '') AS ua_address,
	COALESCE(a.address, '') AS a_address,
	aa.address AS aa_address,
	aa.created_at as created_at
FROM
	emcd.autopay_addresses aa
LEFT JOIN
	emcd.users_accounts ua ON aa.user_account_id = ua.id
LEFT JOIN
	emcd.addresses a ON a.user_account_id = ua.id and token_id is null
WHERE aa.address <> '' and aa.address is not null
AND ua.account_type_id in (1,2)
AND ((ua.address <> '' and ua.address is not null and ua.address <> aa.address) or (a.address <> '' and a.address is not null and a.address <> aa.address))
AND aa.created_at >= $1 order by created_at asc limit $2;
`

const sqlQueryCount = `SELECT
	count(*)
FROM
	emcd.autopay_addresses aa
LEFT JOIN
	emcd.users_accounts ua ON aa.user_account_id = ua.id
LEFT JOIN
	emcd.addresses a ON a.user_account_id = ua.id and token_id is null
WHERE aa.address <> '' and aa.address is not null
AND ua.account_type_id in (1,2)
AND ((ua.address <> '' and ua.address is not null and ua.address <> aa.address) or (a.address <> '' and a.address is not null and a.address <> aa.address))
AND aa.created_at >= $1
`

func (r *migrationRepositoryImpl) GetAddressPersonalMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.AddressPersonalMigrations, error) {
	var totalCount uint64
	if err := r.pgxTransactorOld.Runner(ctx).QueryRow(ctx, sqlQueryCount, lastAtGt).Scan(&totalCount); err != nil {

		return nil, nil, fmt.Errorf("failed scan count: %w", err)
	} else if rows, err := r.pgxTransactorOld.Runner(ctx).Query(ctx, sqlQuery, lastAtGt, limit); err != nil {

		return nil, nil, fmt.Errorf("failed query rows: %w", err)
	} else {
		defer rows.Close()

		var addressMigrationList model_migration.AddressPersonalMigrations
		for rows.Next() {
			var addressMigration model_migration.AddressPersonalMigration
			if err := rows.Scan(
				&addressMigration.UserUuid,
				&addressMigration.CoinCode,
				&addressMigration.MinPayout,
				&addressMigration.UaAddress,
				&addressMigration.AAddress,
				&addressMigration.AaAddress,
				&addressMigration.CreatedAt,
			); err != nil {

				return nil, nil, fmt.Errorf("failed scan rows: %w", err)
			}

			addressMigrationList = append(addressMigrationList, &addressMigration)

		}

		return &totalCount, addressMigrationList, nil
	}
}
