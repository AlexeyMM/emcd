package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressPersonalFilterSql struct {
	*model.AddressPersonalFilter
}

func newAddressPersonalFilterSql(filter *model.AddressPersonalFilter) *addressPersonalFilterSql {

	return &addressPersonalFilterSql{
		AddressPersonalFilter: filter,
	}
}

func (filter *addressPersonalFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.Id != nil {
		query = query.Where(squirrel.Eq{"ap.id": *filter.Id})

	}

	if filter.Address != nil {
		query = query.Where(squirrel.Eq{"ap.address": *filter.Address})

	}

	if filter.UserUuid != nil {
		query = query.Where(squirrel.Eq{"ap.user_uuid": *filter.UserUuid})

	}

	if filter.Network != nil {
		query = query.Where(squirrel.Eq{"ap.network": filter.Network.ToString()})

	}

	if filter.IsDeleted != nil {
		if *filter.IsDeleted {
			query = query.Where(squirrel.NotEq{"ap.deleted_at": nil})

		} else {
			query = query.Where(squirrel.Eq{"ap.deleted_at": nil})

		}
	}

	return query
}
