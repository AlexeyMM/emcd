package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type paginationSql struct {
	*model.Pagination
}

func newPaginationSql(pagination *model.Pagination) *paginationSql {

	return &paginationSql{
		Pagination: pagination,
	}
}

func (p *paginationSql) applyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	query = query.Limit(p.Limit)
	query = query.Offset(p.Offset)

	return query
}
