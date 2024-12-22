package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"code.emcdtech.com/emcd/service/accounting/model"
)

type transactionCollectorFilterSql struct {
	*model.TransactionCollectorFilter
}

func newTransactionCollectorFilterSql(filter *model.TransactionCollectorFilter) *transactionCollectorFilterSql {

	return &transactionCollectorFilterSql{
		TransactionCollectorFilter: filter,
	}
}

func (filter *transactionCollectorFilterSql) ApplyToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.CoinId != nil {
		query = query.Where(squirrel.Eq{"t.coin_id": *filter.CoinId})

	}

	if len(filter.Types) > 0 {
		query = query.Where(squirrel.Eq{"t.type": pq.Array(filter.Types)})

	}

	if filter.CreatedAtGt != nil {
		query = query.Where(squirrel.Gt{"t.created_at": *filter.CreatedAtGt})

	}

	if filter.CreatedAtLte != nil {
		query = query.Where(squirrel.LtOrEq{"t.created_at": *filter.CreatedAtLte})

	}

	return query
}
