package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errTransferNotFound = fmt.Errorf("transfer not found")

type Transfer struct {
	transactor.PgxTransactor
}

func NewTransfer(db *pgxpool.Pool) *Transfer {
	return &Transfer{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (t *Transfer) Add(ctx context.Context, transfer *model.InternalTransfer) error {
	query := squirrel.
		Insert("swap.internal_transfers").
		Columns(
			"id",
			"coin",
			"amount",
			"from_account_id",
			"to_account_id",
			"from_account_type",
			"to_account_type",
			"status",
			"updated_at").
		Values(
			transfer.ID,
			transfer.Coin,
			transfer.Amount,
			transfer.FromAccountID,
			transfer.ToAccountID,
			transfer.FromAccountType,
			transfer.ToAccountType,
			transfer.Status,
			transfer.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = t.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (t *Transfer) Find(ctx context.Context, filter *model.InternalTransferFilter) (model.InternalTransfers, error) {
	return t.find(ctx, filter)
}

func (t *Transfer) FindOne(ctx context.Context, filter *model.InternalTransferFilter) (*model.InternalTransfer, error) {
	transfers, err := t.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(transfers) != 1 {
		if len(transfers) == 0 {
			return nil, errTransferNotFound
		}

		return nil, fmt.Errorf("unexpected number of transfers: %d", len(transfers))
	}
	return transfers[0], nil
}

func (t *Transfer) find(ctx context.Context, filter *model.InternalTransferFilter) (model.InternalTransfers, error) {
	query := squirrel.
		Select(
			"id",
			"coin",
			"amount",
			"from_account_id",
			"to_account_id",
			"from_account_type",
			"to_account_type",
			"status",
			"updated_at").
		From("swap.internal_transfers").
		PlaceholderFormat(squirrel.Dollar)

	query = newInternalTransferFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := t.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var transfers model.InternalTransfers
	for rows.Next() {
		var transfer model.InternalTransfer
		err = rows.Scan(
			&transfer.ID,
			&transfer.Coin,
			&transfer.Amount,
			&transfer.FromAccountID,
			&transfer.ToAccountID,
			&transfer.FromAccountType,
			&transfer.ToAccountType,
			&transfer.Status,
			&transfer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		transfers = append(transfers, &transfer)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return transfers, nil
}

func (t *Transfer) Update(ctx context.Context, internalTransfer *model.InternalTransfer, filter *model.InternalTransferFilter,
	partial *model.InternalTransferPartial) error {
	query := squirrel.
		Update("swap.internal_transfers").
		PlaceholderFormat(squirrel.Dollar)

	query = newInternalTransferFilterSql(filter).applyToUpdateQuery(query)

	query = newInternalTransferPartialSql(partial).applyToSql(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = t.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	internalTransfer.Update(partial)

	return nil
}
