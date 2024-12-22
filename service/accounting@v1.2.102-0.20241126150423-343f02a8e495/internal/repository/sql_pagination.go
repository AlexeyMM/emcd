package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/service/accounting/model"
)

type paginationSql struct {
	*model.Pagination
}

func newPaginationSql(pagination *model.Pagination) *paginationSql {

	return &paginationSql{
		Pagination: pagination,
	}
}

func (p *paginationSql) ApplyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	query = query.Limit(p.Limit)
	query = query.Offset(p.Offset)

	return query
}
