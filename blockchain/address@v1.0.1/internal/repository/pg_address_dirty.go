package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

// CREATE TABLE address_dirty
// (
// address         TEXT      NOT NULL, -- address
// network         TEXT      NOT NULL, -- network
// is_dirty        BOOLEAN   NOT NULL, -- is_dirty
// created_at      TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_dirty_address_network_uniq UNIQUE (address, network)
//
// );

func (r *addressRepositoryImpl) AddOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) error {
	const queryInsert = `insert into address_dirty (address, network, is_dirty, updated_at, created_at)
		values (@address, @network, @is_dirty, @updated_at, @created_at)
		on conflict (address, network) do update 
  SET is_dirty = @is_dirty, 
      updated_at = @updated_at;
		
		`

	addressSql := newAddressDirtySql(address)

	if _, err := r.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}

// func (r *addressRepositoryImpl) UpdateDirtyAddress(ctx context.Context, address *model.AddressDirty, partial *model.AddressDirtyPartial) error {
// 	query := squirrel.
// 		Update("address_dirty").
// 		Where(squirrel.Eq{"address": address.Address}).
// 		Where(squirrel.Eq{"network": address.Network})
//
// 	query = newAddressDirtyPartialSql(partial).applyToQuery(query)
// 	query = query.PlaceholderFormat(squirrel.Dollar)
//
// 	if querySql, args, err := query.ToSql(); err != nil {
// 		return fmt.Errorf("to sql: %w", err)
//
// 	} else if _, err := r.Runner(ctx).Exec(ctx, querySql, args...); err != nil {
//
// 		return fmt.Errorf("exec: %w", err)
// 	} else {
// 		address.Update(partial)
//
// 		return nil
// 	}
// }

func (r *addressRepositoryImpl) GetDirtyAddresses(ctx context.Context, addressFilter *model.AddressDirtyFilter) (model.AddressesDirty, error) {
	query := squirrel.
		Select("ad.address",
			"ad.network",
			"ad.is_dirty",
			"ad.updated_at",
			"ad.created_at",
		).
		From("address_dirty as ad").
		PlaceholderFormat(squirrel.Dollar)

	query = newAddressDirtyFilterSql(addressFilter).applyToQuery(query)
	query = query.Limit(defaultQueryLimit)

	if querySql, args, err := query.ToSql(); err != nil {

		return nil, fmt.Errorf("to sql: %w", err)
	} else if rows, err := r.Runner(ctx).Query(ctx, querySql, args...); err != nil {

		return nil, fmt.Errorf("query rows: %w", err)
	} else {
		defer rows.Close()

		var addressDirtyList model.AddressesDirty
		for rows.Next() {
			var addressDirty model.AddressDirty
			if err := rows.Scan(
				&addressDirty.Address,
				&addressDirty.Network,
				&addressDirty.IsDirty,
				&addressDirty.UpdatedAt,
				&addressDirty.CreatedAt,
			); err != nil {

				return nil, fmt.Errorf("scan rows: %w", err)
			}

			addressDirtyList = append(addressDirtyList, &addressDirty)

		}

		return addressDirtyList, nil
	}
}
