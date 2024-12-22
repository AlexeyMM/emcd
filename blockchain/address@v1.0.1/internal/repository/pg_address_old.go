package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

// CREATE TABLE address_old
// (
// id              UUID      NOT NULL, -- uuid identifier
// address         TEXT      NOT NULL, -- address
// user_uuid       UUID      NOT NULL, -- user (for tables unification)
// address_type    int4      NOT NULL, -- type of address generation
// network         TEXT      NOT NULL, -- network
// user_account_id int4      NOT NULL, -- wallets user account id
// coin            TEXT      NOT NULL, -- coin
// created_at      TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_old_network_user_account_id_idx UNIQUE (network, user_account_id),
// CONSTRAINT address_old_address_idx UNIQUE (address),
//
// PRIMARY KEY (id)
// );

func (r *addressRepositoryImpl) AddOldAddress(ctx context.Context, address *model.AddressOld) error {
	const queryInsert = `insert into address_old (id, address, user_uuid, address_type, network, user_account_id, coin, created_at)
		values (@id, @address, @user_uuid, @address_type, @network, @user_account_id, @coin, @created_at)`

	if address.AddressType.Number() == addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {

		return fmt.Errorf("must using other method for insert derived address") // panic, use other way of add to db
	}

	addressSql := newAddressOldSql(address)

	if _, err := r.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}

func (r *addressRepositoryImpl) GetOldAddresses(ctx context.Context, addressFilter *model.AddressOldFilter) (*uint64, model.AddressesOld, error) {
	var totalCount *uint64
	query := squirrel.
		Select("ao.id",
			"ao.address",
			"ao.user_uuid",
			"ao.address_type",
			"ao.network",
			"ao.user_account_id",
			"ao.coin",
			"ao.created_at",
		).
		From("address_old as ao").
		PlaceholderFormat(squirrel.Dollar)

	query = newAddressOldFilterSql(addressFilter).applyToQuery(query)

	if addressFilter.Pagination != nil {
		var totalCountScan uint64
		queryCount := squirrel.Select("count(*)").
			From("address_old as ao").
			PlaceholderFormat(squirrel.Dollar)

		queryCount = newAddressOldFilterSql(addressFilter).applyToQuery(queryCount)

		if querySql, args, err := queryCount.ToSql(); err != nil {

			return nil, nil, fmt.Errorf("failed to sql: %w", err)
		} else if err := r.PgxTransactor.Runner(ctx).QueryRow(ctx, querySql, args...).Scan(&totalCountScan); err != nil {

			return nil, nil, fmt.Errorf("failed scan count: %w", err)
		} else {
			query = newPaginationSql(addressFilter.Pagination).applyToQuery(query)
			totalCount = &totalCountScan

		}
	} else {
		query = query.Limit(defaultQueryLimit)

	}

	if querySql, args, err := query.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("to sql: %w", err)
	} else if rows, err := r.Runner(ctx).Query(ctx, querySql, args...); err != nil {

		return nil, nil, fmt.Errorf("query rows: %w", err)
	} else {
		defer rows.Close()

		var addressOldList model.AddressesOld
		for rows.Next() {
			var addressOld model.AddressOld
			if err := rows.Scan(
				&addressOld.Id,
				&addressOld.Address,
				&addressOld.UserUuid,
				&addressOld.AddressType,
				&addressOld.Network,
				&addressOld.UserAccountId,
				&addressOld.Coin,
				&addressOld.CreatedAt,
			); err != nil {

				return nil, nil, fmt.Errorf("scan rows: %w", err)
			}

			addressOldList = append(addressOldList, &addressOld)

		}

		return totalCount, addressOldList, nil
	}
}
