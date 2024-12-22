package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressNewPartialSql struct {
	*model.AddressPartial
}

func newAddressNewPartialSql(partial *model.AddressPartial) *addressNewPartialSql {

	return &addressNewPartialSql{
		AddressPartial: partial,
	}
}

func (partial *addressNewPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if partial.Address != nil {
		query = query.Set("address", *partial.Address)

	}

	return query
}
