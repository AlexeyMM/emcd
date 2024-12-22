package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errDepositNotFound = fmt.Errorf("deposit not found")

type Deposit struct {
	transactor.PgxTransactor
}

func NewDeposit(db *pgxpool.Pool) *Deposit {
	return &Deposit{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (d *Deposit) Add(ctx context.Context, deposit *model.Deposit) error {
	sql, args, err := squirrel.
		Insert("swap.deposits").
		Columns(
			"tx_id",
			"swap_id",
			"coin",
			"amount",
			"fee",
			"status",
			"updated_at").
		Values(deposit.TxID,
			deposit.SwapID,
			deposit.Coin,
			deposit.Amount,
			deposit.Fee,
			deposit.Status,
			deposit.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = d.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (d *Deposit) Find(ctx context.Context, filter *model.DepositFilter) (model.Deposits, error) {
	return d.find(ctx, filter)
}

func (d *Deposit) FindOne(ctx context.Context, filter *model.DepositFilter) (*model.Deposit, error) {
	deps, err := d.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(deps) != 1 {
		if len(deps) == 0 {
			return nil, errDepositNotFound
		}

		return nil, fmt.Errorf("unexpected number of deposits: %d", len(deps))
	}

	return deps[0], nil
}

func (d *Deposit) find(ctx context.Context, filter *model.DepositFilter) (model.Deposits, error) {
	query := squirrel.
		Select("tx_id",
			"swap_id",
			"coin",
			"amount",
			"fee",
			"status",
			"updated_at",
		).From("swap.deposits").
		PlaceholderFormat(squirrel.Dollar)

	query = newDepositFilterSql(filter).applyToQuery(query)

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := d.Runner(ctx).Query(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var deposits []*model.Deposit
	for rows.Next() {
		var dep model.Deposit
		err = rows.Scan(&dep.TxID, &dep.SwapID, &dep.Coin, &dep.Amount, &dep.Fee, &dep.Status, &dep.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		deposits = append(deposits, &dep)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return deposits, nil
}
