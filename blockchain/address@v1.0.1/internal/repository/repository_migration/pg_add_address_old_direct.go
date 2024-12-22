package repository_migration

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

func (r *migrationRepositoryImpl) AddOldAddressDirect(ctx context.Context, address *model.AddressOld) error {
	const queryInsert = `insert into address_old (id, address, user_uuid, address_type, network, user_account_id, coin, created_at)
		values (@id, @address, @user_uuid, @address_type, @network, @user_account_id, @coin, @created_at)`

	addressSql := newAddressOldSql(address)

	if _, err := r.pgxTransactorNew.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}
