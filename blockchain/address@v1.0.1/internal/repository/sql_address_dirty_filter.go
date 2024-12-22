package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressDirtyFilterSql struct {
	*model.AddressDirtyFilter
}

func newAddressDirtyFilterSql(filter *model.AddressDirtyFilter) *addressDirtyFilterSql {

	return &addressDirtyFilterSql{
		AddressDirtyFilter: filter,
	}
}

func (filter *addressDirtyFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.Address != nil {
		query = query.Where(squirrel.Eq{"ad.address": *filter.Address})

	}

	if filter.Network != nil {
		query = query.Where(squirrel.Eq{"ad.network": filter.Network.ToString()})

	}

	return query
}
