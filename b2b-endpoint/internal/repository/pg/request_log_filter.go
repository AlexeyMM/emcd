package pg

import (
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
)

type requestLogFilterSql struct {
	*model.RequestLogFilter
}

func newRequestLogFilterSql(filter *model.RequestLogFilter) *requestLogFilterSql {
	return &requestLogFilterSql{
		RequestLogFilter: filter,
	}
}

func (r *requestLogFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if r.ApiKey != nil {
		query = query.Where(squirrel.Eq{"api_key": *r.ApiKey})
	}
	if r.RequestHash != nil {
		query = query.Where(squirrel.Eq{"request_hash": *r.RequestHash})
	}
	return query
}
