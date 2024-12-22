package repository_migration

import (
	"context"
	"fmt"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type migrationRepositoryImpl struct {
	pgxTransactorOld transactor.PgxTransactor
	pgxTransactorNew transactor.PgxTransactor
}

func NewMigrationRepository(poolOld, poolNew *pgxpool.Pool) MigrationRepository {

	return &migrationRepositoryImpl{
		pgxTransactorNew: transactor.NewPgxTransactor(poolNew),
		pgxTransactorOld: transactor.NewPgxTransactor(poolOld),
	}
}

func (r *migrationRepositoryImpl) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {

	return r.pgxTransactorNew.WithinTransaction(ctx, txFn)
}

func (r *migrationRepositoryImpl) AddNewAddressDirect(ctx context.Context, address *model.Address) error {
	const queryInsert = `insert into address (id, address, user_uuid, processing_uuid, address_type, network_group, created_at)
		values (@id, @address, @user_uuid, @processing_uuid, @address_type, @network_group, @created_at)`

	addressSql := newAddressSql(address)

	if _, err := r.pgxTransactorNew.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}
