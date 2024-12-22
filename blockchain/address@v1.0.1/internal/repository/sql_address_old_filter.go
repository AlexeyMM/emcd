package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressOldFilterSql struct {
	*model.AddressOldFilter
}

func newAddressOldFilterSql(filter *model.AddressOldFilter) *addressOldFilterSql {

	return &addressOldFilterSql{
		AddressOldFilter: filter,
	}
}

func (filter *addressOldFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.Id != nil {
		query = query.Where(squirrel.Eq{"ao.id": *filter.Id})

	}

	if filter.Address != nil {
		query = query.Where(squirrel.Eq{"ao.address": *filter.Address})

	}

	if filter.UserUuid != nil {
		query = query.Where(squirrel.Eq{"ao.user_uuid": *filter.UserUuid})

	}

	if filter.AddressType != nil {
		query = query.Where(squirrel.Eq{"ao.address_type": filter.AddressType.Number()})

	}

	if filter.Network != nil {
		query = query.Where(squirrel.Eq{"ao.network": filter.Network.ToString()})

	}

	if filter.UserAccountId != nil {
		query = query.Where(squirrel.Eq{"ao.user_account_id": *filter.UserAccountId})

	}

	if filter.Coin != nil {
		query = query.Where(squirrel.Eq{"ao.coin": *filter.Coin})

	}

	if filter.CreatedAtGt != nil {
		query = query.Where(squirrel.Gt{"ao.created_at": *filter.CreatedAtGt})

	}

	return query
}
