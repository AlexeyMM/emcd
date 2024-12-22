package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressNewFilterSql struct {
	*model.AddressFilter
}

func newAddressNewFilterSql(filter *model.AddressFilter) *addressNewFilterSql {

	return &addressNewFilterSql{
		AddressFilter: filter,
	}
}

func (filter *addressNewFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.Id != nil {
		query = query.Where(squirrel.Eq{"a.id": *filter.Id})

	}

	if filter.Address != nil {
		query = query.Where(squirrel.Eq{"a.address": *filter.Address})

	}

	if filter.UserUuid != nil {
		query = query.Where(squirrel.Eq{"a.user_uuid": *filter.UserUuid})

	}

	if filter.IsProcessing != nil {
		if *filter.IsProcessing {
			query = query.Where("a.processing_uuid <> a.user_uuid")

		} else {
			query = query.Where("a.processing_uuid = a.user_uuid")

		}
	}

	if filter.AddressType != nil {
		query = query.Where(squirrel.Eq{"a.address_type": filter.AddressType.Number()})

	}

	if filter.NetworkGroup != nil {
		query = query.Where(squirrel.Eq{"a.network_group": filter.NetworkGroup.ToString()})

	}

	if filter.CreatedAtGt != nil {
		query = query.Where(squirrel.Gt{"a.created_at": *filter.CreatedAtGt})

	}

	return query
}
