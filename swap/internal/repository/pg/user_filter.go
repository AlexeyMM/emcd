package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type userFilterSql struct {
	*model.UserFilter
}

func newUserFilterSql(filter *model.UserFilter) *userFilterSql {
	return &userFilterSql{
		UserFilter: filter,
	}
}

func (u *userFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if u.ID != nil {
		query = query.Where(squirrel.Eq{"id": *u.ID})
	}
	return query
}
