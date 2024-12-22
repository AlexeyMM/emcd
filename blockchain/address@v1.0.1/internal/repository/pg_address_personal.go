package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

// CREATE TABLE address_personal
// (
// id         UUID      NOT NULL, -- uuid identifier
// address    TEXT      NOT NULL, -- address (possible not uniq)
// user_uuid  UUID      NOT NULL, -- user
// network    TEXT      NOT NULL, -- network group local identifier
// min_payout DOUBLE PRECISION NOT NULL, -- custom minimum payout ge coin.minpay_mining
// deleted_at TIMESTAMP NULL,     -- deleted date
// updated_at TIMESTAMP NOT NULL, -- updated date
// created_at TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_personal_user_uuid_network_group_uniq UNIQUE (user_uuid, network),
//
// PRIMARY KEY (id)
// );

func (r *addressRepositoryImpl) AddPersonalAddress(ctx context.Context, address *model.AddressPersonal) error {
	const queryInsert = `insert into address_personal (id, address, user_uuid, network, min_payout, deleted_at, updated_at, created_at)
		values (@id, @address, @user_uuid, @network, @min_payout, @deleted_at, @updated_at, @created_at)`

	addressSql := newAddressPersonalSql(address)

	if _, err := r.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}

func (r *addressRepositoryImpl) UpdatePersonalAddress(ctx context.Context, address *model.AddressPersonal, addressPartial *model.AddressPersonalPartial) error {
	query := squirrel.
		Update("address_personal").
		Where(squirrel.Eq{"id": address.Id})

	query = newAddressPersonalPartialSql(addressPartial).applyToQuery(query)
	query = query.PlaceholderFormat(squirrel.Dollar)

	if querySql, args, err := query.ToSql(); err != nil {
		return fmt.Errorf("to sql: %w", err)

	} else if _, err := r.Runner(ctx).Exec(ctx, querySql, args...); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {
		address.Update(addressPartial)

		return nil
	}
}

func (r *addressRepositoryImpl) GetPersonalAddresses(ctx context.Context, addressFilter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error) {
	var totalCount *uint64
	query := squirrel.
		Select("ap.id",
			"ap.address",
			"ap.user_uuid",
			"ap.network",
			"ap.min_payout",
			"ap.deleted_at",
			"ap.updated_at",
			"ap.created_at",
		).
		From("address_personal as ap").
		PlaceholderFormat(squirrel.Dollar)

	query = newAddressPersonalFilterSql(addressFilter).applyToQuery(query)

	if addressFilter.Pagination != nil {
		var totalCountScan uint64
		queryCount := squirrel.Select("count(*)").
			From("address_personal as ap").
			PlaceholderFormat(squirrel.Dollar)

		queryCount = newAddressPersonalFilterSql(addressFilter).applyToQuery(queryCount)

		if querySql, args, err := queryCount.ToSql(); err != nil {

			return nil, nil, fmt.Errorf("failed to sql: %w", err)
		} else if err := r.PgxTransactor.Runner(ctx).QueryRow(ctx, querySql, args...).Scan(&totalCountScan); err != nil {

			return nil, nil, fmt.Errorf("failed scan count: %w", err)
		} else {
			query = newPaginationSql(addressFilter.Pagination).applyToQuery(query)
			totalCount = &totalCountScan

		}
	}

	if querySql, args, err := query.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("to sql: %w", err)
	} else if rows, err := r.Runner(ctx).Query(ctx, querySql, args...); err != nil {

		return nil, nil, fmt.Errorf("query rows: %w", err)
	} else {
		defer rows.Close()

		var addressPersonalList model.AddressesPersonal
		for rows.Next() {
			var address model.AddressPersonal
			if err := rows.Scan(
				&address.Id,
				&address.Address,
				&address.UserUuid,
				&address.Network,
				&address.MinPayout,
				&address.DeletedAt,
				&address.UpdatedAt,
				&address.CreatedAt,
			); err != nil {

				return nil, nil, fmt.Errorf("scan rows: %w", err)
			}

			addressPersonalList = append(addressPersonalList, &address)

		}

		return totalCount, addressPersonalList, nil
	}
}
