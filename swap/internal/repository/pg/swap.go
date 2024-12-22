package pg

import (
	"context"
	"errors"
	"fmt"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/b2b/swap/model"
)

var errSwapNotFound = errors.New("swap not found")

type Swap struct {
	transactor.PgxTransactor
}

func NewSwap(db *pgxpool.Pool) *Swap {
	return &Swap{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (s *Swap) Add(ctx context.Context, swap *model.Swap) error {
	query := squirrel.
		Insert("swap.swaps").
		Columns(
			"id",
			"user_id",
			"account_from",
			"coin_from",
			"address_from",
			"network_from",
			"coin_to",
			"address_to",
			"network_to",
			"tag_to",
			"amount_from",
			"tag_from",
			"amount_to",
			"status",
			"start_time",
			"end_time",
			"partner_id").
		Values(
			swap.ID,
			swap.UserID,
			swap.AccountFrom,
			swap.CoinFrom,
			swap.AddressFrom,
			swap.NetworkFrom,
			swap.CoinTo,
			swap.AddressTo,
			swap.NetworkTo,
			swap.TagTo,
			swap.AmountFrom,
			swap.TagFrom,
			swap.AmountTo,
			swap.Status,
			swap.StartTime.UTC(),
			swap.EndTime.UTC(),
			swap.PartnerID,
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = s.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (s *Swap) Find(ctx context.Context, filter *model.SwapFilter) (model.Swaps, error) {
	return s.find(ctx, filter)
}

func (s *Swap) FindOne(ctx context.Context, filter *model.SwapFilter) (*model.Swap, error) {
	swaps, err := s.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(swaps) != 1 {
		if len(swaps) == 0 {
			return nil, errSwapNotFound
		}
		return nil, fmt.Errorf("unexpected nubmer of swaps %d", len(swaps))
	}
	return swaps[0], nil
}

func (s *Swap) CountTotalWithFilter(ctx context.Context, filter *model.SwapFilter) (int, error) {
	query := squirrel.Select("count(1)").
		From("swap.swaps").
		PlaceholderFormat(squirrel.Dollar)

	filterCopy := *filter
	filterCopy.Limit = nil
	filterCopy.Offset = nil
	query = newSwapFilterSql(&filterCopy).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("toSql: %w", err)
	}

	row := s.Runner(ctx).QueryRow(ctx, sql, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}

	return total, nil
}

func (s *Swap) find(ctx context.Context, filter *model.SwapFilter) (model.Swaps, error) {
	query := squirrel.Select(
		"swaps.id",
		"swaps.user_id",
		"swaps.account_from",
		"swaps.coin_from",
		"swaps.address_from",
		"swaps.network_from",
		"swaps.coin_to",
		"swaps.address_to",
		"swaps.network_to",
		"swaps.tag_to",
		"swaps.amount_from",
		"swaps.amount_to",
		"swaps.status",
		"swaps.start_time",
		"swaps.end_time",
		"swaps.tag_from",
		"partner_id").
		From("swap.swaps").
		PlaceholderFormat(squirrel.Dollar)

	query = newSwapFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := s.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var swaps []*model.Swap

	for rows.Next() {
		var swap model.Swap
		err = rows.Scan(
			&swap.ID,
			&swap.UserID,
			&swap.AccountFrom,
			&swap.CoinFrom,
			&swap.AddressFrom,
			&swap.NetworkFrom,
			&swap.CoinTo,
			&swap.AddressTo,
			&swap.NetworkTo,
			&swap.TagTo,
			&swap.AmountFrom,
			&swap.AmountTo,
			&swap.Status,
			&swap.StartTime,
			&swap.EndTime,
			&swap.TagFrom,
			&swap.PartnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		swaps = append(swaps, &swap)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return swaps, nil
}

func (s *Swap) Update(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial) error {
	query := squirrel.
		Update("swap.swaps").
		PlaceholderFormat(squirrel.Dollar)

	query = newSwapFilterSql(filter).applyToUpdateQuery(query)

	query = newSwapPartialSql(partial).applyToQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = s.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	swap.Update(partial)

	return nil
}

func (s *Swap) CountSwapsByStatus(ctx context.Context) (map[model.Status]int, error) {
	query := squirrel.Select(
		"status",
		"count(1)").
		From("swap.swaps").
		GroupBy("status")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}
	rows, err := s.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	m := make(map[model.Status]int)
	for rows.Next() {
		var (
			status model.Status
			count  int
		)
		if err = rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		m[status] = count
	}

	return m, nil
}
