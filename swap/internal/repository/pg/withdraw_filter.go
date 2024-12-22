package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type withdrawFilterSql struct {
	*model.WithdrawFilter
}

func newWithdrawFilterSql(filter *model.WithdrawFilter) *withdrawFilterSql {
	return &withdrawFilterSql{
		WithdrawFilter: filter,
	}
}

func (w *withdrawFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if w.ID != nil {
		query = query.Where(squirrel.Eq{"id": *w.ID})
	}
	if w.SwapID != nil {
		query = query.Where(squirrel.Eq{"swap_id": *w.SwapID})
	}
	return query
}

func (w *withdrawFilterSql) applyToUpdateQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if w.ID != nil {
		query = query.Where(squirrel.Eq{"id": *w.ID})
	}
	if w.SwapID != nil {
		query = query.Where(squirrel.Eq{"swap_id": *w.SwapID})
	}
	return query
}
