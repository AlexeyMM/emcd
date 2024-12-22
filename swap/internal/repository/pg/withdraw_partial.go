package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type withdrawPartialSql struct {
	*model.WithdrawPartial
}

func newWithdrawPartialSql(filter *model.WithdrawPartial) *withdrawPartialSql {
	return &withdrawPartialSql{
		WithdrawPartial: filter,
	}
}

func (s *withdrawPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if s.ID != nil {
		query = query.Set("id", *s.ID)
	}
	if s.Amount != nil {
		query = query.Set("amount", *s.Amount)
	}
	if s.Status != nil {
		query = query.Set("status", *s.Status)
	}
	if s.HashID != nil {
		query = query.Set("hash_id", *s.HashID)
	}
	return query
}
