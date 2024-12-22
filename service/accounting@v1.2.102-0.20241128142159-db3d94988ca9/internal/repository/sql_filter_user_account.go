package repository

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/emcd/service/accounting/model"
)

type userAccountFilterSql struct {
	*model.UserAccountFilter
}

func newUserAccountFilterSql(filter *model.UserAccountFilter) *userAccountFilterSql {

	return &userAccountFilterSql{
		UserAccountFilter: filter,
	}
}

func (filter *userAccountFilterSql) ApplyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.ID != nil {
		query = query.Where(squirrel.Eq{"ua.id": *filter.ID})

	}

	if filter.UserID != nil {
		query = query.Where(squirrel.Eq{"ua.user_id": *filter.UserID})

	}

	if filter.AccountTypeID != nil {
		query = query.Where(squirrel.Eq{"ua.account_type_id": *filter.AccountTypeID})

	}

	if filter.UserIDNew != nil {
		query = query.Where(squirrel.Eq{"ua.user_id_new": *filter.UserIDNew})

	}

	if filter.CoinNew != nil {
		query = query.Where(squirrel.Eq{"ua.coin_new": *filter.CoinNew})

	}

	if filter.IsActive != nil {
		query = query.Where(squirrel.Eq{"ua.is_active": *filter.IsActive})

	}

	if filter.UserIDNewIsNull != nil {
		query = query.Where(squirrel.Eq{"ua.user_id_new": nil})

	}

	if filter.CoinNewIsNull != nil {
		query = query.Where(squirrel.Eq{"ua.coin_new": nil})

	}

	return query
}
