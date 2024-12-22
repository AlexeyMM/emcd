package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type internalTransferFilterSql struct {
	*model.InternalTransferFilter
}

func newInternalTransferFilterSql(filter *model.InternalTransferFilter) *internalTransferFilterSql {
	return &internalTransferFilterSql{
		InternalTransferFilter: filter,
	}
}

func (i *internalTransferFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if i.ID != nil {
		query = query.Where(squirrel.Eq{"id": *i.ID})
	}
	if i.FromAccountID != nil {
		query = query.Where(squirrel.Eq{"from_account_id": *i.FromAccountID})
	}
	if i.IsLast != nil && *i.IsLast {
		query = query.OrderBy("updated_at desc")
		query = query.Limit(1)
	}
	return query
}

func (i *internalTransferFilterSql) applyToUpdateQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if i.ID != nil {
		query = query.Where(squirrel.Eq{"id": *i.ID})
	}
	if i.FromAccountID != nil {
		query = query.Where(squirrel.Eq{"from_account_id": *i.FromAccountID})
	}
	return query
}
