package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type orderFilterSql struct {
	*model.OrderFilter
}

func newOrderFilterSql(filter *model.OrderFilter) *orderFilterSql {
	return &orderFilterSql{
		OrderFilter: filter,
	}
}

func (o *orderFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if o.ID != nil {
		query = query.Where(squirrel.Eq{"id": *o.ID})
	}
	if o.AccountID != nil {
		query = query.Where(squirrel.Eq{"account_id": *o.AccountID})
	}
	if o.IsFirst != nil {
		query = query.Where(squirrel.Eq{"is_first": *o.IsFirst})
	}
	if o.LtStatus != nil {
		query = query.Where(squirrel.Lt{"status": *o.LtStatus})
	}
	return query
}

func (o *orderFilterSql) applyToUpdateQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if o.ID != nil {
		query = query.Where(squirrel.Eq{"id": *o.ID})
	}
	if o.AccountID != nil {
		query = query.Where(squirrel.Eq{"account_id": *o.AccountID})
	}
	if o.IsFirst != nil {
		query = query.Where(squirrel.Eq{"is_first": *o.IsFirst})
	}
	if o.LtStatus != nil {
		query = query.Where(squirrel.Lt{"status": *o.LtStatus})
	}
	return query
}
