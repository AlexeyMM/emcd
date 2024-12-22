package repository_migration

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

// CREATE TABLE migration_date
// (
// table_name TEXT      NOT NULL,
// last_at    TIMESTAMP NOT NULL, -- last migration date
// CONSTRAINT migration_date_uniq UNIQUE (table_name)
// );

func (r *migrationRepositoryImpl) GetMigrationLastAt(ctx context.Context, tableName model_migration.MigrationTableName) (*time.Time, error) {
	query := squirrel.
		Select("m.last_at").
		From("migration_date as m").
		Where(squirrel.Eq{"m.table_name": tableName}).
		PlaceholderFormat(squirrel.Dollar)

	var lastAt time.Time
	if querySql, args, err := query.ToSql(); err != nil {

		return nil, fmt.Errorf("to sql: %w", err)
	} else if err := r.pgxTransactorNew.Runner(ctx).QueryRow(ctx, querySql, args...).Scan(&lastAt); err != nil {

		return nil, fmt.Errorf("scan: %w", err)
	} else {

		return &lastAt, nil
	}
}

func (r *migrationRepositoryImpl) UpdateMigrationLastAt(ctx context.Context, tableName model_migration.MigrationTableName, lastAt time.Time) error {
	query := squirrel.
		Update("migration_date as m").
		Where(squirrel.Eq{"m.table_name": tableName}).
		Set("last_at", lastAt).
		PlaceholderFormat(squirrel.Dollar)

	if querySql, args, err := query.ToSql(); err != nil {
		return fmt.Errorf("to sql: %w", err)

	} else if _, err := r.pgxTransactorNew.Runner(ctx).Exec(ctx, querySql, args...); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}
