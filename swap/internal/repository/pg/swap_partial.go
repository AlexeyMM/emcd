package pg

import (
	"github.com/Masterminds/squirrel"

	"code.emcdtech.com/b2b/swap/model"
)

type swapPartialSql struct {
	*model.SwapPartial
}

func newSwapPartialSql(partial *model.SwapPartial) *swapPartialSql {
	return &swapPartialSql{
		SwapPartial: partial,
	}
}

func (s *swapPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if s.UserID != nil {
		query = query.Set("user_id", *s.UserID)
	}
	if s.Status != nil {
		query = query.Set("status", *s.Status)
		query = query.Where(squirrel.Lt{"status": *s.Status}) // Можем только увеличивать статус. Мера предосторожности
	}
	if s.AmountFrom != nil {
		query = query.Set("amount_from", *s.AmountFrom)
	}
	if s.AmountTo != nil {
		query = query.Set("amount_to", *s.AmountTo)
	}
	if s.StartTime != nil {
		query = query.Set("start_time", *s.StartTime)
	}
	if s.EndTime != nil {
		query = query.Set("end_time", *s.EndTime)
	}
	if s.AddressTo != nil {
		query = query.Set("address_to", *s.AddressTo)
	}
	if s.TagTo != nil {
		query = query.Set("tag_to", *s.TagTo)
	}
	return query
}
