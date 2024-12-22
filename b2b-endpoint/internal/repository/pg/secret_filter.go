package pg

import (
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
)

type secretFilterSql struct {
	*model.SecretFilter
}

func newSecretFilterSql(filter *model.SecretFilter) *secretFilterSql {
	return &secretFilterSql{
		SecretFilter: filter,
	}
}

func (s *secretFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if s.ApiKey != nil {
		query = query.Where(squirrel.Eq{"api_key": *s.ApiKey})
	}
	if s.ClientID != nil {
		query = query.Where(squirrel.Eq{"client_id": *s.ClientID})
	}
	if s.IsActive != nil {
		query = query.Where(squirrel.Eq{"is_active": *s.IsActive})
	}
	return query
}

func (s *secretFilterSql) applyToUpdateQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if s.ApiKey != nil {
		query = query.Where(squirrel.Eq{"api_key": *s.ApiKey})
	}
	if s.ClientID != nil {
		query = query.Where(squirrel.Eq{"client_id": *s.ClientID})
	}
	if s.IsActive != nil {
		query = query.Where(squirrel.Eq{"is_active": *s.IsActive})
	}
	return query
}
