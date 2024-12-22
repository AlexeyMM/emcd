package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/service/accounting/model"
)

type userAccountPartialSql struct {
	*model.UserAccountPartial
}

func newUserAccountPartialSql(partial *model.UserAccountPartial) *userAccountPartialSql {

	return &userAccountPartialSql{
		UserAccountPartial: partial,
	}
}

func (partial *userAccountPartialSql) ApplyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if partial.UserIDNew != nil {
		query = query.Set("user_id_new", *partial.UserIDNew)

	}

	if partial.CoinNew != nil {
		query = query.Set("coin_new", *partial.CoinNew)

	}

	return query
}
