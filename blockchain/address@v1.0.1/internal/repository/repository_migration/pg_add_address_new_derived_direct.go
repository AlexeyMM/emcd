package repository_migration

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

func (r *migrationRepositoryImpl) AddNewDerivedAddressDirect(ctx context.Context, address *model.AddressDerived) error {
	const queryInsert = `
		insert into address_derived (address_uuid, network_group, master_key_id, derived_offset)
		values (@address_uuid, @network_group, @master_key_id, @derived_offset)`

	addressDerivedSql := newAddressDerivedSql(address)

	if _, err := r.pgxTransactorNew.Runner(ctx).Exec(ctx, queryInsert, addressDerivedSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}
