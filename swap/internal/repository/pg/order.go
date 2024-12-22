package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/sdk/log"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errOrderNotFound = errors.New("order not found")

type Order struct {
	transactor.PgxTransactor
}

func NewOrder(db *pgxpool.Pool) *Order {
	return &Order{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (o *Order) Add(ctx context.Context, order *model.Order) error {
	sql, args, err := squirrel.
		Insert("swap.orders").
		Columns("id",
			"swap_id",
			"account_id",
			"category",
			"symbol",
			"direction",
			"amount_from",
			"amount_to",
			"status",
			"is_first").
		Values(
			order.ID,
			order.SwapID,
			order.AccountID,
			order.Category,
			order.Symbol,
			order.Direction,
			order.AmountFrom,
			order.AmountTo,
			order.Status,
			order.IsFirst).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = o.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (o *Order) Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error) {
	return o.find(ctx, filter)
}

func (o *Order) FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error) {
	orders, err := o.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(orders) != 1 {
		if len(orders) == 0 {
			return nil, errOrderNotFound
		}

		return nil, fmt.Errorf("unexpected number of orders: %d", len(orders))
	}
	return orders[0], nil
}

func (o *Order) find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error) {
	query := squirrel.
		Select(
			"id",
			"swap_id",
			"account_id",
			"category",
			"symbol",
			"direction",
			"amount_from",
			"amount_to",
			"status",
			"is_first").
		From("swap.orders").
		PlaceholderFormat(squirrel.Dollar)

	query = newOrderFilterSql(filter).applyToSelectQuery(query)

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := o.
		Runner(ctx).Query(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var orders model.Orders
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.ID,
			&order.SwapID,
			&order.AccountID,
			&order.Category,
			&order.Symbol,
			&order.Direction,
			&order.AmountFrom,
			&order.AmountTo,
			&order.Status,
			&order.IsFirst)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return orders, nil
}

func (o *Order) Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error {
	query := squirrel.
		Update("swap.orders").
		PlaceholderFormat(squirrel.Dollar)

	query = newOrderFilterSql(filter).applyToUpdateQuery(query)
	query = newOrderPartialSql(partial).applyToQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	log.Debug(ctx, "update order: %+v", order)

	_, err = o.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	order.Update(partial)

	return nil
}
