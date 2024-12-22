package repository

import (
	"context"
	"fmt"
	"time"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

const maxRepeat = 10
const repeatTimeout = 100 * time.Millisecond
const defaultQueryLimit = 1000

type addressRepositoryImpl struct {
	transactor.PgxTransactor
}

func NewAddressRepository(pool *pgxpool.Pool) AddressRepository {

	return &addressRepositoryImpl{
		PgxTransactor: transactor.NewPgxTransactor(pool),
	}
}

func (r *addressRepositoryImpl) AddNewCommonAddress(ctx context.Context, address *model.Address) error {
	const queryInsert = `insert into address (id, address, user_uuid, processing_uuid, address_type, network_group, created_at)
		values (@id, @address, @user_uuid, @processing_uuid, @address_type, @network_group, @created_at)`

	if address.AddressType.Number() == addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {

		return fmt.Errorf("must using other method for insert derived address") // panic, use other way of add to db
	}

	addressSql := newAddressNewSql(address)

	if _, err := r.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)
	} else {

		return nil
	}
}

func (r *addressRepositoryImpl) AddNewDerivedAddress(ctx context.Context, address *model.Address, masterKeyId uint32, derivedFunc DerivedFunc) error {
	const queryInsert = `insert into address (id, address, user_uuid, processing_uuid, address_type, network_group, created_at)
		values (@id, @address, @user_uuid, @processing_uuid, @address_type, @network_group, @created_at)`

	if address.AddressType.Number() != addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {

		return fmt.Errorf("must using other method for insert non derived address") // panic, use other way of add to db
	}

	address.Address = fmt.Sprintf("fake:%s", uuid.New()) // non blocking unique fake constraint, derived address must update after derive
	addressSql := newAddressNewSql(address)
	if _, err := r.Runner(ctx).Exec(ctx, queryInsert, addressSql.toNamedArgs()); err != nil {

		return fmt.Errorf("exec: %w", err)

	} else if derivedOffset, err := r.addSeriesAddressDerivedRepeat(ctx, address.Id, masterKeyId, address.NetworkGroup.NetworkGroupEnum, maxRepeat); err != nil {

		return fmt.Errorf("serieses derived: %w", err)
	} else if addressStr, err := derivedFunc(derivedOffset); err != nil {

		return fmt.Errorf("derived func: %w", err)
	} else if err := r.updateNewAddress(ctx, address, &model.AddressPartial{Address: &addressStr}); err != nil {

		return fmt.Errorf("update address: %w", err)
	} else {
		addressDerived := &model.AddressDerived{
			AddressUuid:   address.Id,
			MasterKeyId:   masterKeyId,
			DerivedOffset: derivedOffset,
			NetworkGroup:  address.NetworkGroup,
		}

		address.SetAddressDerived(addressDerived)

		return nil
	}
}

func (r *addressRepositoryImpl) updateNewAddress(ctx context.Context, address *model.Address, addressPartial *model.AddressPartial) error {
	query := squirrel.
		Update("address").
		Where(squirrel.Eq{"id": address.Id})

	query = newAddressNewPartialSql(addressPartial).applyToQuery(query)
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

func (r *addressRepositoryImpl) GetNewAddresses(ctx context.Context, addressFilter *model.AddressFilter) (*uint64, model.Addresses, error) {
	var totalCount *uint64
	query := squirrel.
		Select("a.id",
			"a.address",
			"a.user_uuid",
			"a.processing_uuid",
			"a.address_type",
			"a.network_group",
			"a.created_at",
			"d.master_key_id",
			"d.derived_offset",
		).
		From("address as a").
		LeftJoin("address_derived as d ON a.id = d.address_uuid").
		PlaceholderFormat(squirrel.Dollar)

	query = newAddressNewFilterSql(addressFilter).applyToQuery(query)

	if addressFilter.Pagination != nil {
		var totalCountScan uint64
		queryCount := squirrel.Select("count(*)").
			From("address as a").
			PlaceholderFormat(squirrel.Dollar)

		queryCount = newAddressNewFilterSql(addressFilter).applyToQuery(queryCount)

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

		var addressList model.Addresses
		for rows.Next() {
			var address model.Address
			var masterKeyId *uint32
			var derivedOffset *uint32
			if err := rows.Scan(
				&address.Id,
				&address.Address,
				&address.UserUuid,
				&address.ProcessingUuid,
				&address.AddressType,
				&address.NetworkGroup,
				&address.CreatedAt,
				&masterKeyId,
				&derivedOffset,
			); err != nil {

				return nil, nil, fmt.Errorf("scan rows: %w", err)
			}

			if address.AddressType.AddressType.Number() == addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {
				if masterKeyId == nil {

					return nil, nil, fmt.Errorf("master key cannot be nil for address type: %v", address.AddressType)
				} else if derivedOffset == nil {

					return nil, nil, fmt.Errorf("derived offset cannot be nil for address type: %v", address.AddressType)
				} else {
					addressDerived := &model.AddressDerived{
						AddressUuid:   address.Id,
						MasterKeyId:   *masterKeyId,
						DerivedOffset: *derivedOffset,
						NetworkGroup:  address.NetworkGroup,
					}

					address.SetAddressDerived(addressDerived)
				}
			}

			addressList = append(addressList, &address)

		}

		return totalCount, addressList, nil
	}
}
