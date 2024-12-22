package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type accountFilterSql struct {
	*model.AccountFilter
}

func newAccountFilterSql(filter *model.AccountFilter) *accountFilterSql {
	return &accountFilterSql{
		AccountFilter: filter,
	}
}

func (a *accountFilterSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if a.ID != nil {
		query = query.Where(squirrel.Eq{"a.id": *a.ID})
	}
	return query
}
