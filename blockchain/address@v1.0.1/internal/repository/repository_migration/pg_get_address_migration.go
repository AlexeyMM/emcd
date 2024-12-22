package repository_migration

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

// CREATE TABLE emcd.addresses (
// id serial4 NOT NULL,
// user_account_id int4 NULL,
// coin_id int4 NOT NULL,
// token_id int4 NULL,
// address varchar(128) NOT NULL,
// created_at timestamp NOT NULL,
// deleted_at timestamp NULL,
// address_offset int4 NULL,
// );

func (r *migrationRepositoryImpl) GetAddressMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.AddressMigrations, error) {
	queryCount := squirrel.
		Select("count(*) as count").
		From("emcd.addresses as a").
		Where(squirrel.NotEq{"a.address": ""}).
		Where(squirrel.Gt{"a.created_at": lastAtGt}).
		Where(squirrel.Eq{"a.deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar)

	query := squirrel.
		Select("a.id",
			"a.user_account_id",
			"a.coin_id",
			"a.token_id",
			"a.address",
			"a.network_id",
			"a.created_at",
			"a.deleted_at",
			"a.address_offset",
			"u.new_id",
		).
		From("emcd.addresses as a").
		LeftJoin("emcd.users_accounts ua on ua.id = a.user_account_id").
		LeftJoin("emcd.users u on u.id = ua.user_id").
		Where(squirrel.NotEq{"a.address": ""}).
		Where(squirrel.Gt{"a.created_at": lastAtGt}).
		Where(squirrel.Eq{"a.deleted_at": nil}).
		OrderBy("a.created_at ASC").
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar)

	var totalCount uint64
	if querySqlCount, args, err := queryCount.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("failed to sql count: %w", err)
	} else if err := r.pgxTransactorOld.Runner(ctx).QueryRow(ctx, querySqlCount, args...).Scan(&totalCount); err != nil {

		return nil, nil, fmt.Errorf("failed scan count: %w", err)
	} else if querySql, args, err := query.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("failed to sql: %w", err)
	} else if rows, err := r.pgxTransactorOld.Runner(ctx).Query(ctx, querySql, args...); err != nil {

		return nil, nil, fmt.Errorf("failed query rows: %w", err)
	} else {
		defer rows.Close()

		var addressMigrationList model_migration.AddressMigrations
		for rows.Next() {
			var addressMigration model_migration.AddressMigration
			if err := rows.Scan(&addressMigration.Id,
				&addressMigration.UserAccountId,
				&addressMigration.CoinId,
				&addressMigration.TokenId,
				&addressMigration.Address,
				&addressMigration.NetworkId,
				&addressMigration.CreatedAt,
				&addressMigration.DeletedAt,
				&addressMigration.AddressOffset,
				&addressMigration.UserUuid,
			); err != nil {

				return nil, nil, fmt.Errorf("failed scan rows: %w", err)
			}

			addressMigrationList = append(addressMigrationList, &addressMigration)

		}

		return &totalCount, addressMigrationList, nil
	}
}
