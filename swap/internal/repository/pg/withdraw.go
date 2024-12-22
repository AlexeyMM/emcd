package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errWithdrawNotFound = errors.New("withdraw not found")

type Withdraw struct {
	transactor.PgxTransactor
}

func NewWithdraw(db *pgxpool.Pool) *Withdraw {
	return &Withdraw{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (w *Withdraw) Add(ctx context.Context, withdraw *model.Withdraw) error {
	query := squirrel.Insert("swap.withdraws").
		Columns(
			"id",
			"internal_id",
			"swap_id",
			"hash_id",
			"coin",
			"network",
			"address",
			"tag",
			"amount",
			"include_fee_in_amount",
			"status",
			"explorer_link").
		Values(
			withdraw.ID,
			withdraw.InternalID,
			withdraw.SwapID,
			withdraw.HashID,
			withdraw.Coin,
			withdraw.Network,
			withdraw.Address,
			withdraw.Tag,
			withdraw.Amount,
			withdraw.IncludeFeeInAmount,
			withdraw.Status,
			withdraw.ExplorerLink).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = w.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (w *Withdraw) Find(ctx context.Context, filter *model.WithdrawFilter) (model.Withdraws, error) {
	return w.find(ctx, filter)
}

func (w *Withdraw) FindOne(ctx context.Context, filter *model.WithdrawFilter) (*model.Withdraw, error) {
	ws, err := w.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(ws) != 1 {
		if len(ws) == 0 {
			return nil, errWithdrawNotFound
		}
		return nil, fmt.Errorf("unexpected number of witdraws: %d", len(ws))
	}
	return ws[0], nil
}

func (w *Withdraw) find(ctx context.Context, filter *model.WithdrawFilter) (model.Withdraws, error) {
	query := squirrel.
		Select(
			"id",
			"internal_id",
			"swap_id",
			"hash_id",
			"coin",
			"network",
			"address",
			"tag",
			"amount",
			"include_fee_in_amount",
			"status",
			"explorer_link").
		From("swap.withdraws").
		PlaceholderFormat(squirrel.Dollar)

	query = newWithdrawFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := w.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	defer rows.Close()

	var wothdraws model.Withdraws
	for rows.Next() {
		var withdraw model.Withdraw
		err = rows.Scan(
			&withdraw.ID,
			&withdraw.InternalID,
			&withdraw.SwapID,
			&withdraw.HashID,
			&withdraw.Coin,
			&withdraw.Network,
			&withdraw.Address,
			&withdraw.Tag,
			&withdraw.Amount,
			&withdraw.IncludeFeeInAmount,
			&withdraw.Status,
			&withdraw.ExplorerLink,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		wothdraws = append(wothdraws, &withdraw)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return wothdraws, nil
}

func (w *Withdraw) Update(ctx context.Context, withdraw *model.Withdraw, filter *model.WithdrawFilter, partial *model.WithdrawPartial) error {
	query := squirrel.
		Update("swap.withdraws").
		PlaceholderFormat(squirrel.Dollar)

	query = newWithdrawFilterSql(filter).applyToUpdateQuery(query)

	query = newWithdrawPartialSql(partial).applyToQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = w.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	withdraw.Update(partial)

	return nil
}
