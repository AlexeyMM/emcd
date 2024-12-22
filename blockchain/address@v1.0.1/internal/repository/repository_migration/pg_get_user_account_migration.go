package repository_migration

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

func (r *migrationRepositoryImpl) GetUserAccountMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.UserAccountMigrations, error) {
	queryCount := squirrel.
		Select("count(*) as count").
		From("emcd.users_accounts as ua").
		Where(squirrel.GtOrEq{"ua.created_at": lastAtGt}).
		Where(squirrel.NotEq{"ua.is_active": nil}).
		Where(squirrel.Eq{"ua.is_active": true}).
		Where(squirrel.NotEq{"address": nil}).
		Where(squirrel.NotEq{"address": ""}).
		Where(squirrel.Eq{"ua.account_type_id": enum.WalletAccountTypeID.ToInt32()}).
		PlaceholderFormat(squirrel.Dollar)

	query := squirrel.
		Select("ua.id",
			"ua.user_id",
			"ua.coin_id",
			"ua.account_type_id",
			"ua.address",
			"ua.is_active",
			"ua.created_at",
			"u.new_id",
		).
		From("emcd.users_accounts as ua").
		LeftJoin("emcd.users u on u.id = ua.user_id").
		Where(squirrel.GtOrEq{"ua.created_at": lastAtGt}).
		Where(squirrel.NotEq{"ua.is_active": nil}).
		Where(squirrel.Eq{"ua.is_active": true}).
		Where(squirrel.NotEq{"address": nil}).
		Where(squirrel.NotEq{"address": ""}).
		Where(squirrel.Eq{"ua.account_type_id": enum.WalletAccountTypeID.ToInt32()}).
		OrderBy("ua.created_at ASC").
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

		var userAccountMigrationList model_migration.UserAccountMigrations
		for rows.Next() {
			var userAccountMigration model_migration.UserAccountMigration
			if err := rows.Scan(&userAccountMigration.Id,
				&userAccountMigration.UserId,
				&userAccountMigration.CoinId,
				&userAccountMigration.AccountTypeId,
				&userAccountMigration.Address,
				&userAccountMigration.IsActive,
				&userAccountMigration.CreatedAt,
				&userAccountMigration.UserUuid,
			); err != nil {

				return nil, nil, fmt.Errorf("failed scan rows: %w", err)
			}

			userAccountMigrationList = append(userAccountMigrationList, &userAccountMigration)

		}

		return &totalCount, userAccountMigrationList, nil
	}
}
