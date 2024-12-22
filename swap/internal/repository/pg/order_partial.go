package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type orderPartialSql struct {
	*model.OrderPartial
}

func newOrderPartialSql(partial *model.OrderPartial) *orderPartialSql {
	return &orderPartialSql{
		OrderPartial: partial,
	}
}

func (o *orderPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if o.AmountFrom != nil {
		query = query.Set("amount_from", o.AmountFrom)
	}
	if o.AmountTo != nil {
		query = query.Set("amount_to", o.AmountTo)
	}
	if o.Status != nil {
		query = query.Set("status", o.Status)
	}
	return query
}
