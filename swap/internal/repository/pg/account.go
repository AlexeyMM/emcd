package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errAccountNotFound = fmt.Errorf("account not found")

type Account struct {
	transactor.PgxTransactor
}

func NewAccount(db *pgxpool.Pool) *Account {
	return &Account{
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (a *Account) Add(ctx context.Context, account *model.Account) error {
	queryInsertAccountSql, argsAccount, err := squirrel.
		Insert("swap.accounts").
		Columns("id", "is_valid").
		Values(account.ID, account.IsValid).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql 1: %w", err)
	}

	queryInsertSecretsSql, argsSecrets, err := squirrel.
		Insert("swap.secrets").
		Columns("account_id", "api_key", "api_secret").
		Values(account.ID, account.Keys.ApiKey, account.Keys.ApiSecret).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql 2: %w", err)
	}

	err = a.PgxTransactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err := a.PgxTransactor.Runner(ctx).Exec(ctx, queryInsertAccountSql, argsAccount...)
		if err != nil {
			return fmt.Errorf("queryInsertAccount exec: %w", err)
		}

		account.Keys.AccountID = account.ID
		_, err = a.PgxTransactor.Runner(ctx).Exec(ctx, queryInsertSecretsSql, argsSecrets...)
		if err != nil {
			return fmt.Errorf("queryInsertSecrets exec: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("withinTransaction %w", err)
	}
	return nil
}

func (a *Account) Find(ctx context.Context, filter *model.AccountFilter) (model.Accounts, error) {
	return a.find(ctx, filter)
}

func (a *Account) FindOne(ctx context.Context, filter *model.AccountFilter) (*model.Account, error) {
	accounts, err := a.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(accounts) != 1 {
		if len(accounts) == 0 {
			return nil, errAccountNotFound
		}

		return nil, fmt.Errorf("unexpected number of accounts: %d", len(accounts))
	}

	return accounts[0], nil
}

func (a *Account) find(ctx context.Context, filter *model.AccountFilter) (model.Accounts, error) {
	query := squirrel.
		Select("a.id",
			"a.is_valid",
			"s.api_key",
			"s.api_secret",
		).From("swap.accounts as a").
		Join("swap.secrets as s ON s.account_id = a.id").
		PlaceholderFormat(squirrel.Dollar)

	query = newAccountFilterSql(filter).applyToQuery(query)

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := a.PgxTransactor.Runner(ctx).Query(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	accounts := make(model.Accounts, 0)
	for rows.Next() {
		var (
			account model.Account
			secrets model.Secrets
		)
		account.Keys = &secrets
		err = rows.Scan(&account.ID, &account.IsValid, &secrets.ApiKey, &secrets.ApiSecret)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		secrets.AccountID = account.ID
		accounts = append(accounts, &account)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return accounts, nil
}
