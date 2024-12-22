package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type swapStatusHistoryFilter struct {
	*model.SwapStatusHistoryFilter
}

func newSwapStatusHistoryFilter(filter *model.SwapStatusHistoryFilter) *swapStatusHistoryFilter {
	return &swapStatusHistoryFilter{
		SwapStatusHistoryFilter: filter,
	}
}

func (s *swapStatusHistoryFilter) applyToSelectQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if s.SwapID != nil {
		query = query.Where(squirrel.Eq{"swap_id": *s.SwapID})
	}

	return query
}
