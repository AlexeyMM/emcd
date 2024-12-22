package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type swapFilterSql struct {
	*model.SwapFilter
}

func newSwapFilterSql(filter *model.SwapFilter) *swapFilterSql {
	return &swapFilterSql{
		SwapFilter: filter,
	}
}

func (s *swapFilterSql) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if s.ID != nil {
		query = query.Where(squirrel.Eq{"id": *s.ID})
	}

	if len(s.NotEqStatus) != 0 {
		for i := range s.NotEqStatus {
			query = query.Where(squirrel.NotEq{"status": s.NotEqStatus[i]})
		}
	}

	if s.TxID != nil {
		query = query.LeftJoin("swap.deposits d ON d.swap_id = swaps.id").Where(squirrel.Eq{"d.tx_id": *s.TxID})
	}

	if s.StartTimeFrom != nil {
		query = query.Where(squirrel.GtOrEq{"start_time": *s.StartTimeFrom})
	}

	if s.StartTimeTo != nil {
		query = query.Where(squirrel.LtOrEq{"start_time": *s.StartTimeTo})
	}

	if s.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *s.UserID})
	}

	if s.Email != nil {
		query = query.LeftJoin("swap.users u ON u.id = swaps.user_id").Where(squirrel.Eq{"u.email": *s.Email})
	}

	if s.Limit != nil {
		query = query.Limit(uint64(*s.Limit))
	}

	if s.Offset != nil {
		query = query.Offset(uint64(*s.Offset))
	}

	return query
}

func (s *swapFilterSql) applyToUpdateQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if s.ID != nil {
		query = query.Where(squirrel.Eq{"id": *s.ID})
	}
	if s.NotEqStatus != nil {
		query = query.Where(squirrel.NotEq{"status": s.NotEqStatus})
	}
	if s.TxID != nil {
		query = query.Where(squirrel.Eq{"d.tx_id": *s.TxID})
	}
	return query
}
