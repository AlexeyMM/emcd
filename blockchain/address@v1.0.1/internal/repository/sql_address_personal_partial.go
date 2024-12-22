package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressPersonalPartialSql struct {
	*model.AddressPersonalPartial
}

func newAddressPersonalPartialSql(partial *model.AddressPersonalPartial) *addressPersonalPartialSql {

	return &addressPersonalPartialSql{
		AddressPersonalPartial: partial,
	}
}

func (partial *addressPersonalPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if partial.Address != nil {
		query = query.Set("address", *partial.Address)

	}

	if partial.MinPayout != nil {
		query = query.Set("min_payout", *partial.MinPayout)

	}

	if partial.DeletedAt != nil {
		query = query.Set("deleted_at", *partial.DeletedAt)

	}

	if partial.UpdatedAt != nil {
		query = query.Set("updated_at", *partial.UpdatedAt)

	}

	return query
}
