package pg

import (
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
)

type secretPartialSql struct {
	*model.SecretPartial
}

func newSecretPartialSql(partial *model.SecretPartial) *secretPartialSql {
	return &secretPartialSql{
		SecretPartial: partial,
	}
}

func (s *secretPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if s.IsActive != nil {
		query = query.Set("is_active", *s.IsActive)
	}
	if s.LastUsed != nil {
		lastUsed := *s.LastUsed
		query = query.Set("last_used", lastUsed.UTC())
	}
	return query
}
