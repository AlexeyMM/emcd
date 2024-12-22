package repository_migration

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

func (r *migrationRepositoryImpl) AddPersonalAddressDirect(ctx context.Context, address *model.AddressPersonal) error {
	const queryInsert = `insert into address_personal (id, address, user_uuid, network, min_payout, deleted_at, updated_at, created_at)
		values (@id, @address, @user_uuid, @network, @min_payout, @deleted_at, @updated_at, @created_at)`

	addressSql := newAddressPersonalSql(address)

	if _, err := r.pgxTransactorNew.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}
