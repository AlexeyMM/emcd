package pg

import (
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
)

type ipFilterSql struct {
	*model.IPFilter
}

func newIPFilterSql(filter *model.IPFilter) *ipFilterSql {
	return &ipFilterSql{
		IPFilter: filter,
	}
}

func (s *ipFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if s.ApiKey != nil {
		query = query.Where(squirrel.Eq{"api_key": *s.ApiKey})
	}
	if s.Address != nil {
		query = query.Where(squirrel.Eq{"ip_address": *s.Address})
	}
	return query
}

func (s *ipFilterSql) applyToDeleteQuery(query squirrel.DeleteBuilder) squirrel.DeleteBuilder {
	if s.ApiKey != nil {
		query = query.Where(squirrel.Eq{"api_key": *s.ApiKey})
	}
	if s.Address != nil {
		query = query.Where(squirrel.Eq{"ip_address": *s.Address})
	}
	return query
}
