package pg

import (
	"context"
	"fmt"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/b2b/swap/model"
)

type SwapStatusHistory struct {
	transactor.PgxTransactor
}

func NewSwapStatusHistory(db *pgxpool.Pool) *SwapStatusHistory {
	return &SwapStatusHistory{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (s *SwapStatusHistory) Find(ctx context.Context, filter *model.SwapStatusHistoryFilter) ([]*model.SwapStatusHistoryItem, error) {
	query := squirrel.Select("status", "set_at").
		From("swap.swap_status_history_items")

	query = newSwapStatusHistoryFilter(filter).applyToSelectQuery(query)

	query = query.OrderBy("set_at desc").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := s.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var items []*model.SwapStatusHistoryItem
	for rows.Next() {
		var item model.SwapStatusHistoryItem
		if err := rows.Scan(&item.Status, &item.SetAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (s *SwapStatusHistory) Add(ctx context.Context, swapID uuid.UUID, item *model.SwapStatusHistoryItem) error {
	query := squirrel.Insert("swap.swap_status_history_items").
		Columns("swap_id", "status", "set_at").
		Values(swapID, item.Status, item.SetAt.UTC()).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	if _, err := s.Runner(ctx).Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}
