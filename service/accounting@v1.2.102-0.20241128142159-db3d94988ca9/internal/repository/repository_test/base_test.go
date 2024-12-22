package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type Action struct {
	UUID              string
	Status            string
	CreatedAt         time.Time
	Amount            decimal.Decimal
	Fee               decimal.Decimal
	SenderAccountID   int64
	ReceiverAccountID int64
	CoinID            int64
	Type              int
	ReceiverAddress   string
	Comment           string
	TokenID           int64
	UserID            int64
	SendMessage       bool
}

func TestNullStringSuccessNotNullValue(t *testing.T) {
	value := "hello"
	result := repository.NullString(value)
	require.Equal(t, sql.NullString{String: value, Valid: true}, result)
}

func TestNullStringSuccessNullValue(t *testing.T) {
	value := ""
	result := repository.NullString(value)
	require.Equal(t, sql.NullString{String: value, Valid: false}, result)
}

func TestNullInt64StringSuccessNotNullValue(t *testing.T) {
	value := int64(1)
	result := repository.NullInt64(value)
	require.Equal(t, sql.NullInt64{Int64: value, Valid: true}, result)
}

func TestNullInt64StringSuccessZeroValue(t *testing.T) {
	value := int64(0)
	result := repository.NullInt64(value)
	require.Equal(t, sql.NullInt64{Int64: value, Valid: false}, result)
}

func TestDefaultTimeNowSuccessNotNullValue(t *testing.T) {
	value := time.Now()
	result := repository.DefaultTimeNow(value)
	require.Equal(t, value, result)
}

func TestDefaultTimeNowSuccessNullValue(t *testing.T) {
	value := time.Time{}
	result := repository.DefaultTimeNow(value)
	require.False(t, result.IsZero())
}

func truncateTables(ctx context.Context, t *testing.T, tables ...string) {
	for _, table := range tables {
		_, err := dbPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s cascade", table))
		require.NoError(t, err)
	}
}

func writeCoin(ctx context.Context, pool *pgxpool.Pool, t *testing.T, coin model.Coin) {
	query := `insert into emcd.coins (id, name, description, code, rate) values ($1,$2,$3,$4,$5);`
	_, err := pool.Exec(ctx, query, coin.ID, coin.Name, coin.Description, coin.Code, coin.Rate)
	require.NoError(t, err)
}

func writeUser(ctx context.Context, pool *pgxpool.Pool, t *testing.T, user model.User) {
	query := `insert into emcd.users (id, new_id, username, password) values ($1,$2,$3,$4);`
	_, err := pool.Exec(ctx, query, user.ID, user.NewID, user.Username, user.Password)
	require.NoError(t, err)
}

func writeUserWithNewID(ctx context.Context, pool *pgxpool.Pool, t *testing.T, user model.User) {
	query := `insert into emcd.users (id, username, password,new_id,ref_id) values ($1,$2,$3,$4,$5);`
	_, err := pool.Exec(ctx, query, user.ID, user.Username, user.Password, user.NewID, user.RefID)
	require.NoError(t, err)
}

func writeAccountType(ctx context.Context, pool *pgxpool.Pool, t *testing.T, accountTypeId enum.AccountTypeIdWrapper) {
	query := `insert into emcd.account_types (id, name) values ($1,$2);`
	_, err := pool.Exec(ctx, query, accountTypeId.ToInt(), accountTypeId.ToString())
	require.NoError(t, err)
}

func writeUsersAccounts(ctx context.Context, pool *pgxpool.Pool, t *testing.T, ua model.UserAccount) {
	query := `insert into emcd.users_accounts (id, user_id, coin_id, account_type_id, minpay, is_active) values ($1,$2,$3,$4,$5,$6);`
	_, err := pool.Exec(ctx, query, ua.ID, ua.UserID, ua.CoinID, ua.AccountTypeID.ToInt32(), ua.Minpay, ua.IsActive)
	require.NoError(t, err)
}

func getTransaction(ctx context.Context, pool *pgxpool.Pool, t *testing.T, id int) model.Transaction {
	row := pool.QueryRow(ctx, `SELECT type, created_at, sender_account_id, receiver_account_id, coin_id, amount,
					comment, hash, receiver_address, action_id FROM emcd.transactions where id = $1`, id)
	var tx model.Transaction
	err := row.Scan(&tx.Type, &tx.CreatedAt, &tx.SenderAccountID, &tx.ReceiverAccountID,
		&tx.CoinID, &tx.Amount, &tx.Comment, &tx.Hash, &tx.ReceiverAddress, &tx.ActionID)
	require.NoError(t, err)

	return tx
}

func getOperations(ctx context.Context, pool *pgxpool.Pool, t *testing.T, transactionID int) []*model.Operation {
	rows, err := pool.Query(ctx, `SELECT type, transaction_id, account_id, coin_id, created_at, amount FROM
				emcd.operations where transaction_id = $1`, transactionID)
	require.NoError(t, err)
	defer rows.Close()

	result := make([]*model.Operation, 0)
	for rows.Next() {
		var o model.Operation
		err = rows.Scan(&o.Type, &o.TransactionID, &o.AccountID, &o.CoinID, &o.CreatedAt, &o.Amount)
		require.NoError(t, err)
		result = append(result, &o)
	}

	return result
}

func getBlock(ctx context.Context, pool *pgxpool.Pool, t *testing.T) model.Block {
	row := pool.QueryRow(ctx, `SELECT id, block_transaction_id, unblock_transaction_id, unblock_to_account_id,
				blocked_till FROM emcd.transactions_blocks`)
	var b model.Block
	var unlockTransactionId sql.NullInt64
	err := row.Scan(&b.ID, &b.BlockTransactionID, &unlockTransactionId, &b.UnblockToAccountID, &b.BlockedTill)
	require.NoError(t, err)

	if unlockTransactionId.Valid {
		b.UnblockTransactionID = unlockTransactionId.Int64
	}
	return b
}

func writeAction(ctx context.Context, pool *pgxpool.Pool, t *testing.T, action Action) {
	query := `insert into emcd.wallets_accounting_actions (action_id,status,created_at,amount,fee,sender_account_id,receiver_account_id,coin_id,type,receiver_address,comment,token_id,user_id,send_message)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);`
	_, err := pool.Exec(ctx, query, action.UUID, action.Status, action.CreatedAt, action.Amount, action.Fee,
		action.SenderAccountID, action.ReceiverAccountID, action.CoinID, action.Type, action.ReceiverAddress, action.Comment, action.TokenID, action.UserID, action.SendMessage)
	require.NoError(t, err)
}
