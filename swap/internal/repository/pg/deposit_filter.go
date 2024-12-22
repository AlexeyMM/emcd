package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type depositFilterSql struct {
	*model.DepositFilter
}

func newDepositFilterSql(filter *model.DepositFilter) *depositFilterSql {
	return &depositFilterSql{
		DepositFilter: filter,
	}
}

func (d *depositFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if d.SwapID != nil {
		query = query.Where(squirrel.Eq{"swap_id": *d.SwapID})
	}
	if d.UpdatedAt != nil {
		query = query.Where(squirrel.Gt{"updated_at": *d.UpdatedAt})
	}
	return query
}
