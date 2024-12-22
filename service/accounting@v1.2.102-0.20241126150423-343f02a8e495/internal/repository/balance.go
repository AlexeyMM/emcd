package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/service/accounting/model/enum"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/service/accounting/model"
)

//go:generate mockery --name Balance

type Balance interface {
	View(ctx context.Context, userID, coinID int64, accountTypeID enum.AccountTypeId, totalBalance bool) (decimal.Decimal, error)
	Change(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) (transactionID int, err error)
	ChangeWithBlock(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) (transactionID int, err error)
	ChangeWithUnblock(ctx context.Context, sqlTx pgx.Tx, unblockTransaction *model.Transaction) (transactionID int, err error)
	FindOperations(ctx context.Context, userID, coinID int) ([]*model.OperationSelectionWithBlock, error)
	FindBatchOperations(ctx context.Context, usersWithCoins map[int]int) (map[int][]*model.OperationSelection, error)
	FindTransactions(ctx context.Context, types []int, userID, accountTypeID int, coinIDs []int, from time.Time) ([]*model.Transaction, error)
	FindTransactionsByCollectorFilter(ctx context.Context, filter *model.TransactionCollectorFilter) (*uint64, []*model.Transaction, error)
	GetTransactionsByActionID(ctx context.Context, actionID string) ([]*model.Transaction, error)
	FindTransactionsWithBlocks(ctx context.Context, blockedTill time.Time) ([]*model.TransactionSelectionWithBlock, error)
	GetTransactionByIDs(ctx context.Context, sqlTx pgx.Tx, ids ...int64) ([]*model.Transaction, error)
	GetTransactionByID(ctx context.Context, id int) (*model.TransactionSelectionWithBlock, error)
	FindLastBlockTimeBalances(ctx context.Context, userAccountIDs []int64) (map[int]decimal.Decimal, error)
	FindBalancesDiffMining(ctx context.Context, users []*model.UserBeforePayoutMining) (map[int]decimal.Decimal, error)
	FindBalancesDiffWallet(ctx context.Context, users []*model.UserBeforePayoutWallet) ([]model.UserWalletDiff, error)
	ChangeMultiple(ctx context.Context, trs []*model.Transaction) error
	GetBalances(ctx context.Context, userID int32, walletCoinsIDs []int, walletCoinsStrIDs map[int]string) ([]*model.Balance, error)
	GetBalanceByCoin(ctx context.Context, userID int32, coinID int) (*model.Balance, error)
	GetPaid(ctx context.Context, userID int32, coinID int, from, to time.Time) (decimal.Decimal, error)
	GetCoinsSummary(ctx context.Context, userID int32, walletCoinsIDs []int, walletCoinsStrIDs map[int]string) ([]*model.CoinSummary, error)

	// These methods are for `magic` API:
	GetAccountTypeIDByID(ctx context.Context, sqlTx pgx.Tx, id int) (enum.AccountTypeId, error)
	GetBlockAccountIDFor31Type(ctx context.Context, sqlTx pgx.Tx, coinID int, username string) (int32, error)
	GetBlockAccountIDFor57Type(ctx context.Context, sqlTx pgx.Tx, coinID int, username string) (int32, error)
	GetBlockAccountIDBySenderAccountID(ctx context.Context, sqlTx pgx.Tx, id int) (int32, error)
	GetTransactionIDByAction(ctx context.Context, actionID string, txType int, amount string) (int, error)
	GetUserIDsByAccountID(ctx context.Context, accountID int) (uuid.UUID, int32, int32, error)
	GetUserAccountIDByOldID(ctx context.Context, userID int32, accountType enum.AccountTypeId, coinID int) (int, error)
	GetUserAccountIDByNewID(ctx context.Context, userID uuid.UUID, accountType enum.AccountTypeId, coinID int) (int, error)
	CreateUsersAccount(ctx context.Context, userID int32, accountType enum.AccountTypeId, coinID int, minpay float64) (int, error)
	CreateAccountPool(ctx context.Context, accountID int) error
	CreateAccountReferral(ctx context.Context, accountID, coinID int) error

	FindOperationsAndTransactions(ctx context.Context, queryParams *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error)
	GetBalanceBeforeTransaction(ctx context.Context, accountID, transactionID int64) (decimal.Decimal, error)

	GetP2PAdminId(ctx context.Context) (int, error)
}

type Processing string

const (
	processingTypeDefault = Processing("default")
	processingTypeBlock   = Processing("block")
	processingTypeUnblock = Processing("unblock")
)

type balance struct {
	pool                  *pgxpool.Pool
	whiteListBalanceUsers []string
}

func NewBalance(pool *pgxpool.Pool, whiteListBalanceUsers []string) Balance {
	return &balance{
		pool:                  pool,
		whiteListBalanceUsers: whiteListBalanceUsers,
	}
}

func (r *balance) View(ctx context.Context, userID, coinID int64, accountTypeID enum.AccountTypeId, totalBalance bool) (decimal.Decimal, error) {
	query := `
	SELECT
		(SELECT COALESCE(SUM(o.amount), 0) FROM emcd.operations o WHERE o.account_id = ua.id)
	FROM emcd.users_accounts ua
	WHERE ua.is_active IS TRUE
	  AND ua.user_id = $1
	  AND ua.coin_id = $2
	  AND ua.account_type_id = $3
`
	if totalBalance {
		query = fmt.Sprintf(`
	SELECT unblocked_funds + blocked_funds
	from (SELECT (SELECT COALESCE(SUM(o.amount), 0)
				  FROM emcd.users_accounts ua
						   INNER JOIN emcd.operations o ON ua.id = o.account_id
						   INNER JOIN emcd.transactions_blocks tb ON o.transaction_id = tb.block_transaction_id
				  WHERE ua.user_id = $1
					AND ua.is_active = true
					AND ua.coin_id = $2
					AND ua.account_type_id = $3
					AND tb.unblock_transaction_id IS NULL) AS blocked_funds,
				 (%s)                                      AS unblocked_funds) as funds
	`, query)
	}

	var result decimal.NullDecimal
	row := r.pool.QueryRow(ctx, query, userID, coinID, accountTypeID)
	err := row.Scan(&result)
	if !result.Valid {
		return decimal.Decimal{}, fmt.Errorf("scan funds: %w", pgx.ErrNoRows)
	}
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("scan funds: %w", err)
	}

	return result.Decimal, nil
}

func (r *balance) Change(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) (int, error) {
	transactionID, err := r.createBillingTransactions(ctx, sqlTx, processingTypeDefault, transaction)
	if err != nil {
		return 0, fmt.Errorf("createTransactions: %w", err)
	}

	return transactionID, nil
}

func (r *balance) ChangeWithBlock(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) (int, error) {

	transactionID, err := r.createBillingTransactions(ctx, sqlTx, processingTypeBlock, transaction)
	if err != nil {
		return 0, fmt.Errorf("createBlockTransactions: %w", err)
	}

	return transactionID, nil
}

func (r *balance) ChangeWithUnblock(ctx context.Context, sqlTx pgx.Tx, unblockTransaction *model.Transaction) (int, error) {

	transactionID, err := r.createBillingTransactions(ctx, sqlTx, processingTypeUnblock, unblockTransaction)
	if err != nil {
		return 0, fmt.Errorf("createUnblockTransactions: %w", err)
	}

	return transactionID, nil
}

func (r *balance) FindOperations(ctx context.Context, userID, coinID int) ([]*model.OperationSelectionWithBlock, error) {
	query := `
		SELECT 
			o.amount, 
			o.account_id, 
			o.coin_id, 
			ua.coin_id, 
			ua.account_type_id, 
			ua.is_active, 
			o.type, 
			o.created_at, 
			o.transaction_id,
			tb.unblock_transaction_id, 
			tb.unblock_to_account_id 
		FROM emcd.operations o
				 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
				 LEFT JOIN emcd.transactions_blocks tb ON tb.block_transaction_id = o.transaction_id
		WHERE ua.user_id = $1
		AND ua.coin_id = $2;`
	result := make([]*model.OperationSelectionWithBlock, 0)
	rows, err := r.pool.Query(ctx, query, userID, coinID)
	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			o                    model.OperationSelectionWithBlock
			unblockTransactionId sql.NullInt64
			unblockToAccountID   sql.NullInt64
		)
		err = rows.Scan(
			&o.Amount,
			&o.AccountID,
			&o.OperationCoinID,
			&o.UserAccountCoinID,
			&o.AccountTypeID,
			&o.IsActive,
			&o.Type,
			&o.CreatedAt,
			&o.TransactionID,
			&unblockTransactionId,
			&unblockToAccountID,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}
		o.UnblockTransactionID = unblockTransactionId.Int64
		o.UnblockToAccountID = unblockToAccountID.Int64
		result = append(result, &o)
	}

	return result, nil
}

func (r *balance) FindBatchOperations(ctx context.Context, usersWithCoins map[int]int) (map[int][]*model.OperationSelection, error) {
	query := `
		SELECT $1, o.amount, o.account_id, o.coin_id, ua.coin_id, ua.account_type_id, ua.is_active, o.type,
			o.created_at, o.transaction_id 
		FROM emcd.operations o
				 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
		WHERE ua.user_id = $2
		AND ua.coin_id = $3;`
	batch := pgx.Batch{}
	for userID, coinID := range usersWithCoins {
		batch.Queue(query, strconv.Itoa(userID), userID, coinID)
	}

	batchResult := r.pool.SendBatch(ctx, &batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			sdkLog.Error(ctx, "can't close batch: %v", err)
		}
	}(batchResult)

	result := make(map[int][]*model.OperationSelection, len(usersWithCoins))
	for i := 0; i < len(usersWithCoins); i++ {
		rows, err := batchResult.Query()
		if err != nil {
			return nil, fmt.Errorf("batchResult.Query: %w", err)
		}

		for rows.Next() {
			var (
				o              model.OperationSelection
				userIDAsString string
			)
			err = rows.Scan(&userIDAsString, &o.Amount, &o.AccountID, &o.OperationCoinID, &o.UserAccountCoinID,
				&o.AccountTypeID, &o.IsActive, &o.Type, &o.CreatedAt, &o.TransactionID)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("can't scan: %w", err)
			}

			userID, err := strconv.Atoi(userIDAsString)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("convert to int %w", err)
			}
			result[userID] = append(result[userID], &o)
		}

		rows.Close()
	}

	return result, nil
}

func (r *balance) FindTransactions(ctx context.Context, types []int, userID, accountTypeID int, coinIDs []int, from time.Time) (
	[]*model.Transaction, error) {
	query := `
		SELECT t.id,
			   t.type,
			   t.created_at,
			   t.sender_account_id,
			   t.receiver_account_id,
			   t.coin_id,
			   t.token_id,
			   t.amount,
			   t.comment,
			   t.hash,
			   t.receiver_address,
               t.hashrate,
               t.from_referral_id,
               t.action_id
		FROM emcd.transactions t
				 INNER JOIN emcd.users_accounts ua ON ua.id = t.receiver_account_id
				 INNER JOIN emcd.users u ON u.id = ua.user_id
		WHERE t.type = ANY ($1)
		  AND u.id = $2
		  AND ua.account_type_id = $3
		  AND ua.coin_id = ANY ($4)
		  AND t.created_at > $5
		ORDER BY t.created_at;
		`

	rows, err := r.pool.Query(ctx, query, types, userID, accountTypeID, coinIDs, from)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var result []*model.Transaction
	for rows.Next() {
		var (
			comment           sql.NullString
			hash              sql.NullString
			receiverAddress   sql.NullString
			receiverAccountID sql.NullInt64
			hashrate          sql.NullInt64
			tokenID           sql.NullInt64
			fromReferralId    sql.NullInt64
			amount            decimal.NullDecimal
			actionID          sql.NullString
			t                 model.Transaction
		)
		err = rows.Scan(&t.ID, &t.Type, &t.CreatedAt, &t.SenderAccountID, &receiverAccountID,
			&t.CoinID, &tokenID, &amount, &comment, &hash, &receiverAddress, &hashrate, &fromReferralId, &actionID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		t.Comment = comment.String
		t.Hash = hash.String
		t.Hashrate = hashrate.Int64
		t.ReceiverAddress = receiverAddress.String
		t.ReceiverAccountID = receiverAccountID.Int64
		t.TokenID = tokenID.Int64
		t.Amount = amount.Decimal
		t.FromReferralId = fromReferralId.Int64
		t.ActionID = actionID.String

		result = append(result, &t)
	}

	return result, nil
}

func (r *balance) FindTransactionsByCollectorFilter(ctx context.Context, filter *model.TransactionCollectorFilter) (*uint64, []*model.Transaction, error) {
	var totalCount *uint64
	query := squirrel.
		Select("t.id",
			"t.type",
			"t.created_at",
			"t.sender_account_id",
			"t.receiver_account_id",
			"t.coin_id",
			"t.token_id",
			"t.amount",
			"t.comment",
			"t.hash",
			"t.receiver_address",
			"t.hashrate",
			"t.from_referral_id",
			"t.action_id",
		).From("emcd.transactions t").
		PlaceholderFormat(squirrel.Dollar)

	if filter.Pagination != nil {
		var totalCountScan uint64
		queryCount := squirrel.Select("count(*)").
			From("emcd.transactions t").
			PlaceholderFormat(squirrel.Dollar)

		queryCount = newTransactionCollectorFilterSql(filter).ApplyToQuery(queryCount)

		if querySql, args, err := queryCount.ToSql(); err != nil {

			return nil, nil, fmt.Errorf("failed to sql: %w", err)
		} else if err := r.pool.QueryRow(ctx, querySql, args...).Scan(&totalCountScan); err != nil {

			return nil, nil, fmt.Errorf("failed scan count: %w", err)
		} else {
			query = newPaginationSql(filter.Pagination).ApplyToQuery(query)
			query = query.OrderBy("t.id ASC")
			totalCount = &totalCountScan

		}
	}

	query = newTransactionCollectorFilterSql(filter).ApplyToQuery(query)

	if querySql, args, err := query.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("failed to sql: %w", err)
	} else if rows, err := r.pool.Query(ctx, querySql, args...); err != nil {

		return nil, nil, fmt.Errorf("failed query rows: %w", err)
	} else {
		defer rows.Close()

		var result []*model.Transaction
		for rows.Next() {
			var (
				comment           sql.NullString
				hash              sql.NullString
				receiverAddress   sql.NullString
				receiverAccountID sql.NullInt64
				hashrate          sql.NullInt64
				tokenID           sql.NullInt64
				fromReferralId    sql.NullInt64
				amount            decimal.NullDecimal
				actionID          sql.NullString
				t                 model.Transaction
			)
			err = rows.Scan(&t.ID, &t.Type, &t.CreatedAt, &t.SenderAccountID, &receiverAccountID,
				&t.CoinID, &tokenID, &amount, &comment, &hash, &receiverAddress, &hashrate, &fromReferralId, &actionID)
			if err != nil {
				return totalCount, nil, fmt.Errorf("failed scan: %w", err)
			}
			t.Comment = comment.String
			t.Hash = hash.String
			t.Hashrate = hashrate.Int64
			t.ReceiverAddress = receiverAddress.String
			t.ReceiverAccountID = receiverAccountID.Int64
			t.TokenID = tokenID.Int64
			t.Amount = amount.Decimal
			t.FromReferralId = fromReferralId.Int64
			t.ActionID = actionID.String

			result = append(result, &t)
		}

		return totalCount, result, nil
	}
}

func (r *balance) FindTransactionsWithBlocks(ctx context.Context, blockedTill time.Time) (
	[]*model.TransactionSelectionWithBlock, error) {
	query := `
		SELECT tb.id, t.receiver_account_id, tb.unblock_to_account_id, t.amount, t.coin_id, t.type, t.action_id 
		FROM emcd.transactions_blocks tb
				 INNER JOIN emcd.transactions t ON t.id = tb.block_transaction_id
		WHERE unblock_transaction_id IS NULL
		  AND tb.blocked_till < $1;
	`
	result := make([]*model.TransactionSelectionWithBlock, 0)
	rows, err := r.pool.Query(ctx, query, blockedTill)
	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			receiverAccountID sql.NullInt64
			amount            decimal.NullDecimal
			actionID          sql.NullString
			t                 model.TransactionSelectionWithBlock
		)
		err = rows.Scan(&t.BlockID, &receiverAccountID, &t.UnblockToAccountID, &amount, &t.CoinID, &t.Type, &actionID)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}
		t.ReceiverAccountID = receiverAccountID.Int64
		t.Amount = amount.Decimal
		t.ActionID = actionID.String
		result = append(result, &t)
	}

	return result, nil
}

func (r *balance) GetTransactionByIDs(ctx context.Context, sqlTx pgx.Tx, ids ...int64) ([]*model.Transaction, error) {
	const selectTransactionByIdsSQL = `
SELECT t.id,
	   t.type,
	   t.created_at,
	   t.sender_account_id,
	   t.receiver_account_id,
	   t.coin_id,
	   t.token_id,
	   t.amount,
	   t.comment,
	   t.hash,
	   t.receiver_address,
	   t.hashrate,
	   t.from_referral_id,
	   t.action_id
FROM emcd.transactions t
WHERE t.id  = ANY ($1)
ORDER BY t.created_at;`

	rows, err := sqlTx.Query(ctx, selectTransactionByIdsSQL, ids)
	if err != nil {
		return nil, fmt.Errorf("execute selectTransactionByIdsSQL: %w", err)
	}
	defer rows.Close()

	var result []*model.Transaction
	for rows.Next() {
		var (
			comment           sql.NullString
			hash              sql.NullString
			receiverAddress   sql.NullString
			receiverAccountID sql.NullInt64
			hashrate          sql.NullInt64
			tokenID           sql.NullInt64
			fromReferralId    sql.NullInt64
			amount            decimal.NullDecimal
			actionID          sql.NullString
		)
		t := new(model.Transaction)

		err = rows.Scan(&t.ID, &t.Type, &t.CreatedAt, &t.SenderAccountID, &receiverAccountID,
			&t.CoinID, &tokenID, &amount, &comment, &hash, &receiverAddress, &hashrate, &fromReferralId, &actionID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		t.Comment = comment.String
		t.Hash = hash.String
		t.Hashrate = hashrate.Int64
		t.ReceiverAddress = receiverAddress.String
		t.ReceiverAccountID = receiverAccountID.Int64
		t.TokenID = tokenID.Int64
		t.Amount = amount.Decimal
		t.FromReferralId = fromReferralId.Int64
		t.ActionID = actionID.String

		result = append(result, t)
	}

	return result, nil
}

func (r *balance) GetTransactionByID(ctx context.Context, id int) (*model.TransactionSelectionWithBlock, error) {
	query := `
		SELECT
		    t.sender_account_id, t.receiver_account_id, tb.unblock_to_account_id, t.coin_id, t.amount, t.type,
		    COALESCE(tb.unblock_transaction_id, 0) AS unblock_transaction_id, t.action_id
		FROM
		    emcd.transactions t
		LEFT JOIN
		     emcd.transactions_blocks tb ON tb.block_transaction_id = t.id
		WHERE
		    t.id = $1;
	`
	row := r.pool.QueryRow(ctx, query, id)
	var (
		receiverAccountID sql.NullInt64
		actionID          sql.NullString
		amount            decimal.NullDecimal
		t                 model.TransactionSelectionWithBlock
	)
	err := row.Scan(&t.SenderAccountID, &receiverAccountID, &t.UnblockToAccountID, &t.CoinID,
		&amount, &t.Type, &t.UnblockTransactionID, &actionID)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	t.ReceiverAccountID = receiverAccountID.Int64
	t.Amount = amount.Decimal
	t.ActionID = actionID.String

	return &t, nil
}

func (r *balance) FindLastBlockTimeBalances(ctx context.Context, userAccountIDs []int64) (map[int]decimal.Decimal, error) {
	query := `
		SELECT ua.user_id, COALESCE(SUM(amount), 0)
			FROM emcd.operations o
					 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
			WHERE ua.id = $1
			  AND o.created_at <= (SELECT max(o.created_at)
								   FROM emcd.operations o
											INNER JOIN emcd.transactions block_tr ON o.transaction_id = block_tr.id
											INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
											LEFT JOIN emcd.transactions_blocks tb ON block_tr.id = tb.block_transaction_id
											LEFT JOIN emcd.transactions unblock_tr ON tb.unblock_transaction_id = unblock_tr.id
								   WHERE o.type = 21
									 AND o.account_id = $1
									 AND (tb.id IS NULL OR (unblock_tr.receiver_account_id != block_tr.sender_account_id)))
			GROUP BY ua.user_id;
	`
	batch := pgx.Batch{}
	for _, uaID := range userAccountIDs {
		batch.Queue(query, uaID)
	}

	batchResult := r.pool.SendBatch(ctx, &batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			sdkLog.Error(ctx, "can't close batch: %v", err)
		}
	}(batchResult)

	result := make(map[int]decimal.Decimal, len(userAccountIDs))
	for i := 0; i < len(userAccountIDs); i++ {
		rows, err := batchResult.Query()
		if err != nil {
			return nil, fmt.Errorf("batchResult.Query: %w", err)
		}

		for rows.Next() {
			var (
				userID int
				sum    float64
			)
			err = rows.Scan(&userID, &sum)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan: %w", err)
			}
			result[userID] = decimal.NewFromFloat(sum)
		}
		rows.Close()
	}

	return result, nil
}

func (r *balance) FindBalancesDiffMining(ctx context.Context, users []*model.UserBeforePayoutMining) (map[int]decimal.Decimal, error) {
	query := `
		SELECT $1,
			   (COALESCE(sum(amount), 0) - COALESCE((SELECT amount FROM emcd.transactions WHERE id = $2), 0)) as diff
		FROM emcd.operations o
				 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
		WHERE ua.user_id = $3
		  AND ua.coin_id = $4
		  AND ua.account_type_id = $5
		  AND o.created_at > $6
		  AND o.created_at < (SELECT created_at FROM emcd.transactions WHERE id = $2)
		GROUP BY ua.user_id;
	`
	batch := pgx.Batch{}
	for _, u := range users {
		batch.Queue(query, strconv.Itoa(int(u.BlockID)), u.BlockID, u.UserID, u.CoinID, u.AccountTypeID, u.LastPay)
	}

	batchResult := r.pool.SendBatch(ctx, &batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			sdkLog.Error(ctx, "can't close batch: %v", err)
		}
	}(batchResult)

	result := make(map[int]decimal.Decimal, len(users))
	for i := 0; i < len(users); i++ {
		rows, err := batchResult.Query()
		if err != nil {
			return nil, fmt.Errorf("batchResult.Query: %w", err)
		}

		for rows.Next() {
			var (
				blockIDAsString string
				diff            float64
			)
			err = rows.Scan(&blockIDAsString, &diff)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan: %w", err)
			}
			blockID, err := strconv.Atoi(blockIDAsString)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("convert to int %w", err)
			}
			result[blockID] = decimal.NewFromFloat(diff)
		}

		rows.Close()
	}

	return result, nil
}

func (r *balance) FindBalancesDiffWallet(ctx context.Context, users []*model.UserBeforePayoutWallet) ([]model.UserWalletDiff, error) {
	query := `
		SELECT ua.user_id, $1, 
			   (COALESCE(sum(amount), 0) - COALESCE((SELECT amount FROM emcd.transactions WHERE id = $2), 0)) as diff
		FROM emcd.operations o
				 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
		WHERE ua.user_id = $3 
		  AND ua.coin_id = $4
		  AND ua.account_type_id = $5
		  AND o.created_at < (SELECT created_at FROM emcd.transactions WHERE id = $2)
		GROUP BY ua.user_id;
	`

	entitiesCount := 0
	batch := pgx.Batch{}
	for _, u := range users {
		for _, tx := range u.TransactionIDs {
			entitiesCount++
			batch.Queue(query, strconv.Itoa(int(tx)), tx, u.UserID, u.CoinID, u.AccountTypeID)
		}
	}

	batchResult := r.pool.SendBatch(ctx, &batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			sdkLog.Error(ctx, "can't close batch: %v", err)
		}
	}(batchResult)

	result := make([]model.UserWalletDiff, 0)
	for i := 0; i < entitiesCount; i++ {
		rows, err := batchResult.Query()
		if err != nil {
			return nil, fmt.Errorf("batchResult.Query: %w", err)
		}

		for rows.Next() {
			var (
				userID          int
				blockIDAsString string
				diff            float64
			)
			err = rows.Scan(&userID, &blockIDAsString, &diff)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan: %w", err)
			}
			blockID, err := strconv.Atoi(blockIDAsString)
			if err != nil {
				rows.Close()
				return nil, fmt.Errorf("strconv.Atoi: %w", err)
			}

			result = append(
				result,
				model.UserWalletDiff{UserID: int64(userID), BlockID: int64(blockID), Diff: decimal.NewFromFloat(diff)},
			)
		}

		rows.Close()
	}

	return result, nil
}

func (r *balance) ChangeMultiple(ctx context.Context, trs []*model.Transaction) error {
	sqlTx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin sql tx: %w", err)
	}
	defer func() {
		if err != nil {
			for i := range trs {
				trs[i].ID = 0
			}
			if rollbackErr := sqlTx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				sdkLog.Error(ctx, "can't rollback: %v", rollbackErr)
			}
		}
	}()
	for i := range trs {
		transactionID, err := r.createTransactions(ctx, sqlTx, trs[i])
		if err != nil {
			return fmt.Errorf("createTransactions: %w", err)
		}
		trs[i].ID = int64(transactionID)

		err = r.createOperations(ctx, sqlTx, trs[i])
		if err != nil {
			return fmt.Errorf("createOperations: %w", err)
		}
	}
	err = sqlTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit: %w", err)
	}

	return nil
}

const selectGetBalancesOld = `
	SELECT c.id AS coin_id,
		(SELECT COALESCE(SUM(o0.amount), 0)
			FROM emcd.operations o0
			INNER JOIN emcd.users_accounts ua0 ON ua0.id = o0.account_id
			WHERE o0.account_id = ua0.id
			  AND ua0.user_id = $1
			  AND ua0.coin_id = c.id
			  AND ua0.account_type_id = 1
		) AS wallet_balance,
		(SELECT COALESCE(SUM(o1.amount), 0)
			FROM emcd.operations o1
			INNER JOIN emcd.users_accounts ua1 ON ua1.id = o1.account_id
			WHERE o1.account_id = ua1.id
			  AND ua1.user_id = $1
			  AND ua1.coin_id = c.id
			  AND ua1.account_type_id = 2
		) AS mining_balance,
		(SELECT COALESCE(SUM(o2.amount), 0)
			FROM emcd.operations o2
			INNER JOIN emcd.users_accounts ua2 ON ua2.id = o2.account_id
			WHERE o2.account_id = ua2.id
			  AND ua2.user_id = $1
			  AND ua2.coin_id = c.id
			  AND ua2.account_type_id = 3
			  AND ua2.is_active IS TRUE
		) AS coinholds_balance,
		(SELECT COALESCE(SUM(o3.amount), 0)
			FROM emcd.operations o3
			INNER JOIN emcd.users_accounts ua3 ON ua3.id = o3.account_id
			WHERE o3.account_id = ua3.id
			  AND ua3.user_id = $1
			  AND ua3.coin_id = c.id
			  AND ua3.account_type_id = 7
		) AS p2p_balance,
		(SELECT COALESCE(SUM(o4.amount), 0)
			FROM emcd.operations o4
			INNER JOIN emcd.users_accounts ua4 ON ua4.id = o4.account_id
			INNER JOIN emcd.transactions_blocks tb0 ON tb0.block_transaction_id = o4.transaction_id
			WHERE o4.account_id = ua4.id
			  AND o4.type IN (34, 36)
			  AND ua4.user_id = $1
			  AND ua4.coin_id = c.id
			  AND ua4.account_type_id = 5
			  AND tb0.unblock_transaction_id IS NULL
			  AND tb0.unblock_to_account_id = (SELECT id
												FROM emcd.users_accounts
												WHERE user_id = $1
												  AND account_type_id = 1
												  AND coin_id = c.id)
		) AS blocked_balance_coinhold,
		ABS((SELECT COALESCE(SUM(o5.amount), 0)
			FROM emcd.operations o5
			INNER JOIN emcd.users_accounts ua5 ON ua5.id = o5.account_id
			INNER JOIN emcd.transactions_blocks tb1 ON tb1.block_transaction_id = o5.transaction_id
			WHERE o5.account_id = ua5.id
			  AND o5.type = 31
			  AND ua5.user_id = $1
			  AND ua5.coin_id = c.id
			  AND ua5.account_type_id = 1
			  AND tb1.unblock_transaction_id IS NULL
		)) AS blocked_balance_free_withdraw,
	    ABS((SELECT COALESCE(SUM(o6.amount), 0)
			FROM emcd.operations o6
			INNER JOIN emcd.users_accounts ua6 ON ua6.id = o6.account_id
			INNER JOIN emcd.transactions_blocks tb2 ON tb2.block_transaction_id = o6.transaction_id
			WHERE o6.account_id = ua6.id
			  AND o6.type = 66
			  AND ua6.user_id = $1
			  AND ua6.coin_id = c.id
			  AND ua6.account_type_id = 7
			  AND tb2.unblock_transaction_id IS NULL
		)) AS blocked_balance_p2p
	FROM emcd.coins c
	WHERE c.id = ANY($2)`

func (r *balance) GetBalancesOLD(ctx context.Context, userID int32, walletCoinsIDs []int, walletCoinsStrIDs map[int]string) ([]*model.Balance, error) {
	rows, err := r.pool.Query(ctx, selectGetBalancesOld, userID, walletCoinsIDs)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var balances []*model.Balance
	for rows.Next() {
		b := &model.Balance{
			WalletBalance:               decimal.Decimal{},
			MiningBalance:               decimal.Decimal{},
			CoinholdsBalance:            decimal.Decimal{},
			P2pBalance:                  decimal.Decimal{},
			BlockedBalanceP2p:           decimal.Decimal{},
			BlockedBalanceCoinhold:      decimal.Decimal{},
			BlockedBalanceFreeWithdraw:  decimal.Decimal{},
			BlockedBalanceMiningPayouts: decimal.Decimal{},
			CoinID:                      "",
		}

		coinID := 0
		if err = rows.Scan(
			&coinID,
			&b.WalletBalance,
			&b.MiningBalance,
			&b.CoinholdsBalance,
			&b.P2pBalance,
			&b.BlockedBalanceCoinhold,
			&b.BlockedBalanceFreeWithdraw,
			&b.BlockedBalanceP2p,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		coinName, ok := walletCoinsStrIDs[coinID]
		if !ok {
			return nil, fmt.Errorf("unexpected int coin_id: %d", coinID)
		}
		b.CoinID = coinName
		balances = append(balances, b)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return balances, nil
}

func (r *balance) GetBalanceByCoinOLD(ctx context.Context, userID int32, coinID int) (*model.Balance, error) {
	var (
		b          model.Balance
		coinIDScan int
	)
	err := r.pool.QueryRow(ctx, selectGetBalancesOld, userID, []int{coinID}).Scan(
		&coinIDScan,
		&b.WalletBalance,
		&b.MiningBalance,
		&b.CoinholdsBalance,
		&b.P2pBalance,
		&b.BlockedBalanceCoinhold,
		&b.BlockedBalanceFreeWithdraw,
		&b.BlockedBalanceP2p)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	return &b, nil
}

const selectGetBalancesNew = `
SELECT COALESCE(sum(o.amount),0) as balance, o.coin_id as coin, ua.account_type_id as type FROM emcd.operations o 
INNER JOIN emcd.transactions t on t.id = o.transaction_id
INNER JOIN emcd.users_accounts ua on ua.id = o.account_id
WHERE ua.user_id=$1 and ua.account_type_id IN (1,2,3,7) AND o.coin_id = ANY($2)
GROUP BY ua.account_type_id, o.coin_id;
`
const selectGetBlocked = `
SELECT COALESCE(sum(t.amount),0) as balance, t.coin_id as coin, ua.account_type_id as acc_type, t.type as t_type FROM emcd.transactions t
INNER JOIN emcd.users_accounts ua on ua.id = t.sender_account_id
LEFT JOIN emcd.transactions_blocks tb ON t.id = tb.block_transaction_id
WHERE user_id=$1 and ua.account_type_id IN (1,2,3,7) and tb.id IS NOT NULL and tb.unblock_transaction_id IS NULL AND t.coin_id = ANY($2)
GROUP BY ua.account_type_id, t.coin_id, t.type ORDER BY t.coin_id, ua.account_type_id;
`

const (
	WalletBalance               = "WalletBalance"
	MiningBalance               = "MiningBalance"
	CoinholdsBalance            = "CoinholdsBalance"
	P2pBalance                  = "P2pBalance"
	BlockedBalanceCoinhold      = "BlockedBalanceCoinhold"
	BlockedBalanceFreeWithdraw  = "BlockedBalanceFreeWithdraw"
	BlockedBalanceP2p           = "BlockedBalanceP2p"
	BlockedBalanceMiningPayouts = "BlockedBalanceMiningPayouts"
)

func (r *balance) getBalancesData(ctx context.Context, userID int32, walletCoinsIDs []int) (map[int64]map[string]decimal.Decimal, error) {

	result := make(map[int64]map[string]decimal.Decimal)

	itemNames := [8]string{WalletBalance, MiningBalance, CoinholdsBalance, P2pBalance, BlockedBalanceCoinhold,
		BlockedBalanceFreeWithdraw, BlockedBalanceP2p, BlockedBalanceMiningPayouts}

	for _, coinID := range walletCoinsIDs {
		item := make(map[string]decimal.Decimal)
		for _, itemName := range itemNames {
			item[itemName] = decimal.Decimal{}
		}
		result[int64(coinID)] = item
	}

	rows, err := r.pool.Query(ctx, selectGetBalancesNew, userID, walletCoinsIDs)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("r.pool.Query: %w", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var (
			balance decimal.Decimal
			coinID  int64
			accType enum.AccountTypeIdWrapper
		)

		if err = rows.Scan(
			&balance,
			&coinID,
			&accType,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		switch accType.AccountTypeId {
		case enum.WalletAccountTypeID:
			result[coinID][WalletBalance] = balance
		case enum.MiningAccountTypeID:
			result[coinID][MiningBalance] = balance
		case enum.CoinholdAccountTypeID:
			result[coinID][CoinholdsBalance] = balance
		case enum.P2PUserAccountTypeIDNotUsedBefore:
			result[coinID][P2pBalance] = balance
		default:
			// pass // TODO: is it correct?
		}
	}

	rows2, err := r.pool.Query(ctx, selectGetBlocked, userID, walletCoinsIDs)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("r.pool.Query: %w", err)
		}
	}
	defer rows2.Close()

	for rows2.Next() {
		var (
			balance decimal.Decimal
			coinID  int64
			accType enum.AccountTypeIdWrapper
			trType  int64
		)

		if err = rows2.Scan(
			&balance,
			&coinID,
			&accType,
			&trType,
		); err != nil {
			return nil, fmt.Errorf("rows2.Scan: %w", err)
		}

		const trType21PleaseNameIt = 21
		const trType31PleaseNameIt = 31
		const trType66PleaseNameIt = 66

		switch accType.AccountTypeId {
		case enum.WalletAccountTypeID:
			switch trType {
			case trType31PleaseNameIt:
				result[coinID][BlockedBalanceFreeWithdraw] = balance
			case trType66PleaseNameIt:
				result[coinID][BlockedBalanceP2p] = decimal.Sum(result[coinID][BlockedBalanceP2p], balance)
			}
		case enum.MiningAccountTypeID:
			if trType == trType21PleaseNameIt {
				result[coinID][BlockedBalanceMiningPayouts] = balance
			}
		case enum.CoinholdAccountTypeID:
			result[coinID][BlockedBalanceCoinhold] = decimal.Sum(result[coinID][BlockedBalanceCoinhold], balance)
		case enum.P2PUserAccountTypeIDNotUsedBefore:
			result[coinID][BlockedBalanceP2p] = decimal.Sum(result[coinID][BlockedBalanceP2p], balance)
		default:
			// pass // TODO: is it correct?

		}

	}

	return result, nil
}

func (r *balance) GetBalances(ctx context.Context, userID int32, walletCoinsIDs []int, walletCoinsStrIDs map[int]string) ([]*model.Balance, error) {

	result, err := r.getBalancesData(ctx, userID, walletCoinsIDs)
	if err != nil {
		return nil, fmt.Errorf("r.getBalancesData: %w", err)
	}

	var balances []*model.Balance

	for coinId, item := range result {
		coinName, ok := walletCoinsStrIDs[int(coinId)]
		if !ok {
			return nil, fmt.Errorf("unexpected int coin_id: %d", coinId)
		}
		b := &model.Balance{
			WalletBalance:               item[WalletBalance],
			MiningBalance:               item[MiningBalance],
			CoinholdsBalance:            item[CoinholdsBalance],
			P2pBalance:                  item[P2pBalance],
			BlockedBalanceP2p:           item[BlockedBalanceP2p],
			BlockedBalanceCoinhold:      item[BlockedBalanceCoinhold],
			BlockedBalanceFreeWithdraw:  item[BlockedBalanceFreeWithdraw],
			BlockedBalanceMiningPayouts: item[BlockedBalanceMiningPayouts],
			CoinID:                      coinName,
		}
		balances = append(balances, b)
	}

	return balances, nil
}

func (r *balance) GetBalanceByCoin(ctx context.Context, userID int32, coinID int) (*model.Balance, error) {

	result, err := r.getBalancesData(ctx, userID, []int{coinID})
	if err != nil {
		return nil, fmt.Errorf("getBalancesData: %w", err)
	}

	for coinId, item := range result {
		if coinId == int64(coinID) {
			return &model.Balance{
				WalletBalance:               item[WalletBalance],
				MiningBalance:               item[MiningBalance],
				CoinholdsBalance:            item[CoinholdsBalance],
				P2pBalance:                  item[P2pBalance],
				BlockedBalanceP2p:           item[BlockedBalanceP2p],
				BlockedBalanceCoinhold:      item[BlockedBalanceCoinhold],
				BlockedBalanceFreeWithdraw:  item[BlockedBalanceFreeWithdraw],
				BlockedBalanceMiningPayouts: item[BlockedBalanceMiningPayouts],
				CoinID:                      "", // TODO: why empty?
			}, nil
		}
	}

	return nil, fmt.Errorf("getBalancesData not found by coinId: %d, userId: %d", coinID, userID)
}

const getPaid = `
	SELECT ABS(COALESCE(SUM(CASE WHEN type IN (21, 51) THEN amount ELSE 0 END), 0)) AS paid
			FROM emcd.operations o
			INNER JOIN emcd.users_accounts ua1 ON ua1.id = o.account_id
			WHERE o.account_id = ua1.id
			  AND ua1.user_id = $1
			  AND ua1.coin_id = $2
			  AND ua1.account_type_id = 2
			  AND o.created_at > $3
			  AND o.created_at <= $4;
	`

func (r *balance) GetPaid(ctx context.Context, userID int32, coinID int, from, to time.Time) (decimal.Decimal, error) {
	var paid decimal.Decimal
	err := r.pool.QueryRow(ctx, getPaid, userID, coinID, from, to).Scan(&paid)
	if err != nil {
		return decimal.Zero, fmt.Errorf("queryRow: %w", err)
	}
	return paid, nil
}

const selectGetCoinsSummary = `
	SELECT c.id AS coin_id,
		TRUNC(
			(SELECT COALESCE(SUM(amount), 0)
				FROM emcd.operations o0
				INNER JOIN emcd.users_accounts ua0 ON ua0.id = o0.account_id
				WHERE o0.account_id = ua0.id
				  AND ua0.user_id = $1
				  AND ua0.coin_id = c.id
				  AND ua0.account_type_id = 1
			)::DECIMAL +
			(SELECT COALESCE(SUM(amount), 0)
				FROM emcd.operations o1
				INNER JOIN emcd.users_accounts ua1 ON ua1.id = o1.account_id
				WHERE o1.account_id = ua1.id
				  AND ua1.user_id = $1
				  AND ua1.coin_id = c.id
				  AND ua1.account_type_id = 2
			)::DECIMAL +
			(SELECT COALESCE(SUM(amount), 0)
				FROM emcd.operations o2
				INNER JOIN emcd.users_accounts ua2 ON ua2.id = o2.account_id
				WHERE o2.account_id = ua2.id
				  AND ua2.user_id = $1
				  AND ua2.coin_id = c.id
				  AND ua2.account_type_id = 3
				  AND ua2.is_active IS TRUE
			)::DECIMAL +
			(SELECT COALESCE(SUM(amount), 0)
				FROM emcd.operations o3
				INNER JOIN emcd.users_accounts ua3 ON ua3.id = o3.account_id
				INNER JOIN emcd.transactions_blocks tb0 ON tb0.block_transaction_id = o3.transaction_id
				WHERE o3.account_id = ua3.id
				  AND o3.type IN (34, 36)
				  AND ua3.user_id = $1
				  AND ua3.coin_id = c.id
				  AND ua3.account_type_id = 5
				  AND tb0.unblock_transaction_id IS NULL
				  AND tb0.unblock_to_account_id = (SELECT id
													FROM emcd.users_accounts
													WHERE user_id = $1
													  AND account_type_id = 1
													  AND coin_id = c.id)
			)::DECIMAL +
			(SELECT COALESCE(SUM(amount), 0)
				FROM emcd.operations o4
				INNER JOIN emcd.users_accounts ua4 ON ua4.id = o4.account_id
				WHERE o4.account_id = ua4.id
				  AND ua4.user_id = $1
				  AND ua4.coin_id = c.id
				  AND ua4.account_type_id = 7
			)::DECIMAL +
			ABS((SELECT COALESCE(SUM(o5.amount), 0)
				FROM emcd.operations o5
				INNER JOIN emcd.users_accounts ua5 ON ua5.id = o5.account_id
				INNER JOIN emcd.transactions_blocks tb1 ON tb1.block_transaction_id = o5.transaction_id
				WHERE o5.account_id = ua5.id
				  AND o5.type = 31
				  AND ua5.user_id = $1
				  AND ua5.coin_id = c.id
				  AND ua5.account_type_id = 1
				  AND tb1.unblock_transaction_id IS NULL
			))::DECIMAL
		, 8) AS total_amount
	FROM emcd.coins c
	WHERE c.id = ANY($2)`

func (r *balance) GetCoinsSummary(ctx context.Context, userID int32, walletCoinsIDs []int, walletCoinsStrIDs map[int]string) ([]*model.CoinSummary, error) {
	rows, err := r.pool.Query(ctx, selectGetCoinsSummary, userID, walletCoinsIDs)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var coinsSummary []*model.CoinSummary
	for rows.Next() {
		cs := &model.CoinSummary{
			TotalAmount: decimal.Decimal{},
			CoinID:      "",
		}
		coinID := 0
		if err = rows.Scan(
			&coinID,
			&cs.TotalAmount,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		coinName, ok := walletCoinsStrIDs[coinID]
		if !ok {
			return nil, fmt.Errorf("unexpected int coin_id: %d", coinID)
		}
		cs.CoinID = coinName
		coinsSummary = append(coinsSummary, cs)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return coinsSummary, nil
}

func (r *balance) GetAccountTypeIDByID(ctx context.Context, sqlTx pgx.Tx, id int) (enum.AccountTypeId, error) {
	query := `SELECT account_type_id from emcd.users_accounts where id = $1;`
	var accountTypeID int
	err := sqlTx.QueryRow(ctx, query, id).Scan(&accountTypeID)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return enum.AccountTypeId(accountTypeID), nil
}

func (r *balance) GetBlockAccountIDFor31Type(ctx context.Context, sqlTx pgx.Tx, coinID int, username string) (int32, error) {
	query := `
		SELECT id
		FROM emcd.users_accounts ua
		WHERE ua.user_id = (SELECT id FROM emcd.users WHERE username = $1)
		  AND ua.coin_id = $2
		  AND ua.account_type_id = 5;
	`
	var id int32
	err := sqlTx.QueryRow(ctx, query, username, coinID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return id, nil
}

func (r *balance) GetBlockAccountIDFor57Type(ctx context.Context, sqlTx pgx.Tx, coinID int, username string) (int32, error) {
	query := `
		SELECT block_ua.id
		FROM emcd.users_accounts ua
			INNER JOIN emcd.users u ON u.id = ua.user_id
			LEFT JOIN emcd.users_accounts block_ua
			  ON block_ua.user_id = ua.user_id AND block_ua.coin_id = ua.coin_id AND block_ua.account_type_id = 5
		WHERE ua.account_type_id = 1
		  AND ua.coin_id = $1
		  AND u.username = $2;
	`
	var id int32
	err := sqlTx.QueryRow(ctx, query, coinID, username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return id, nil
}

func (r *balance) GetBlockAccountIDBySenderAccountID(ctx context.Context, sqlTx pgx.Tx, id int) (int32, error) {
	query := `
		SELECT block_ua.id
		FROM emcd.users_accounts ua
			INNER JOIN emcd.users_accounts block_ua
			  ON block_ua.account_type_id = 5 AND block_ua.user_id = ua.user_id AND block_ua.coin_id = ua.coin_id
		WHERE ua.id = $1;
	`
	var blockAccountID int32
	err := sqlTx.QueryRow(ctx, query, id).Scan(&blockAccountID)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return blockAccountID, nil
}

func (r *balance) GetTransactionIDByAction(ctx context.Context, actionID string, txType int, amount string) (int, error) {
	query := `
		SELECT
			id
		FROM
			emcd.transactions
		WHERE
			action_id = $1
			AND type = $2
			AND amount = $3
	`

	var txID int

	err := r.pool.QueryRow(ctx, query, actionID, txType, amount).Scan(&txID)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return txID, nil
}

func (r *balance) GetTransactionsByActionID(ctx context.Context, actionID string) ([]*model.Transaction, error) {
	query := `
		SELECT t.id,
			   t.type,
			   t.created_at,
			   t.sender_account_id,
			   t.receiver_account_id,
			   t.coin_id,
			   t.token_id,
			   t.amount,
			   t.comment,
			   t.hash,
			   t.receiver_address,
               t.hashrate,
               t.from_referral_id,
               t.action_id
		FROM emcd.transactions t
		WHERE t.action_id = $1
		ORDER BY t.created_at;
		`

	rows, err := r.pool.Query(ctx, query, actionID)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var result []*model.Transaction
	for rows.Next() {
		var (
			comment           sql.NullString
			hash              sql.NullString
			receiverAddress   sql.NullString
			receiverAccountID sql.NullInt64
			hashrate          sql.NullInt64
			tokenID           sql.NullInt64
			fromReferralId    sql.NullInt64
			amount            decimal.NullDecimal
			t                 model.Transaction
		)
		err = rows.Scan(&t.ID, &t.Type, &t.CreatedAt, &t.SenderAccountID, &receiverAccountID,
			&t.CoinID, &tokenID, &amount, &comment, &hash, &receiverAddress, &hashrate, &fromReferralId, &t.ActionID)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		t.Comment = comment.String
		t.Hash = hash.String
		t.Hashrate = hashrate.Int64
		t.ReceiverAddress = receiverAddress.String
		t.ReceiverAccountID = receiverAccountID.Int64
		t.TokenID = tokenID.Int64
		t.Amount = amount.Decimal
		t.FromReferralId = fromReferralId.Int64

		result = append(result, &t)
	}

	return result, nil
}

func (r *balance) createTransactions(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) (int, error) {
	query := `
	INSERT INTO emcd.transactions (
			type, 
			created_at, 
			sender_account_id, 
			receiver_account_id, 
			coin_id, 
			amount, 
			comment, 
			hash, 
			receiver_address, 
			token_id,
			action_id,
	        hashrate,                       
			from_refferal_id
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	row := sqlTx.QueryRow(
		ctx,
		query,
		transaction.Type,
		DefaultTimeNow(transaction.CreatedAt),
		transaction.SenderAccountID,
		transaction.ReceiverAccountID,
		transaction.CoinID,
		transaction.Amount,
		NullString(transaction.Comment),
		NullString(transaction.Hash),
		NullString(transaction.ReceiverAddress),
		NullInt64(transaction.TokenID),
		NullString(transaction.ActionID),
		NullInt64(transaction.Hashrate),
		NullInt64(transaction.FromReferralId),
	)
	var transactionID int
	err := row.Scan(&transactionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return transactionID, fmt.Errorf("scan: %w", err)
	}

	return transactionID, nil
}

func (r *balance) createOperations(ctx context.Context, sqlTx pgx.Tx, transaction *model.Transaction) error {
	query := `
	INSERT INTO emcd.operations (type, transaction_id, account_id, coin_id, created_at, amount)
		VALUES ($1, $2, $3, $4, $5, $6), ($1, $2, $7, $4, $5, $8)`

	_, err := sqlTx.Exec(
		ctx,
		query,
		transaction.Type,
		transaction.ID,
		transaction.SenderAccountID,
		transaction.CoinID,
		DefaultTimeNow(transaction.CreatedAt),
		transaction.Amount.Neg(),
		transaction.ReceiverAccountID,
		transaction.Amount,
	)
	if err != nil {
		return fmt.Errorf("can't insert in operations: %w", err)
	}

	return nil
}

func (r *balance) createBillingTransactions(ctx context.Context, sqlTx pgx.Tx, processingType Processing, transaction *model.Transaction) (int, error) {

	query := `
SELECT emcd.accountingCreateTransactions($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);`

	row := sqlTx.QueryRow(
		ctx,
		query,
		pq.Array(r.whiteListBalanceUsers),
		processingType,
		transaction.ActionID,
		transaction.ReceiverAccountID,
		transaction.SenderAccountID,
		transaction.Amount,
		transaction.Type,
		NullString(transaction.Comment),
		NullString(transaction.Hash),
		NullInt64(transaction.Hashrate),
		NullString(transaction.ReceiverAddress),
		NullInt64(transaction.TokenID),
		transaction.BlockedTill,
		transaction.UnblockAccountId,
		transaction.CreatedAt,
		NullInt64(transaction.FromReferralId),
	)
	var transactionID int

	err := row.Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("repo: create tr, type: %d, action_id: %v, error: %w", transaction.Type, transaction.ActionID, err)
	}

	const insertOutboxTransactionSQL = `
INSERT INTO outbox_transactions(transaction_id) VALUES($1);
`
	if _, err := sqlTx.Exec(ctx, insertOutboxTransactionSQL, transactionID); err != nil {
		return 0, fmt.Errorf("execute insertExportTransactionSQL: %w", err)
	}
	return transactionID, nil
}

func (r *balance) GetUserIDsByAccountID(ctx context.Context, accountID int) (uuid.UUID, int32, int32, error) {
	var (
		newID uuid.UUID
		oldID int32
		refID int32
	)
	err := r.pool.QueryRow(ctx, `SELECT u.new_id, u.id, u.ref_id FROM emcd.users_accounts ua JOIN emcd.users u ON ua.user_id=u.id WHERE ua.id=$1`, accountID).
		Scan(&newID, &oldID, &refID)
	if err != nil {
		return uuid.Nil, 0, 0, fmt.Errorf("queryRow: %w", err)
	}
	return newID, oldID, refID, nil
}

func (r *balance) GetUserAccountIDByOldID(ctx context.Context, userID int32, accountType enum.AccountTypeId, coinID int) (int, error) {
	var accountID int
	err := r.pool.QueryRow(ctx, `SELECT id FROM emcd.users_accounts WHERE user_id=$1 AND account_type_id=$2 AND coin_id=$3`,
		userID, accountType, coinID).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("queryRow: %w", err)
	}
	return accountID, nil
}

func (r *balance) GetUserAccountIDByNewID(ctx context.Context, userID uuid.UUID, accountType enum.AccountTypeId, coinID int) (int, error) {
	var accountID int
	err := r.pool.QueryRow(ctx, `SELECT ua.id FROM emcd.users_accounts ua JOIN emcd.users u ON ua.user_id=u.id WHERE u.new_id=$1 AND ua.account_type_id=$2 AND ua.coin_id=$3`,
		userID, accountType, coinID).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("queryRow: %w", err)
	}
	return accountID, nil
}

func (r *balance) CreateUsersAccount(ctx context.Context, userID int32, accountType enum.AccountTypeId, coinID int, minpay float64) (int, error) {
	query := `INSERT INTO emcd.users_accounts (user_id, coin_id, account_type_id, minpay, is_active, created_at, updated_at)
	SELECT $1, $2, $3, $4, TRUE, NOW(), NOW()
	WHERE NOT EXISTS(
			SELECT NULL
			FROM emcd.users_accounts
			WHERE (user_id, coin_id, account_type_id) = ($1, $2, $3)
		)
	RETURNING id`
	var id int
	if err := r.pool.QueryRow(ctx, query, userID, coinID, accountType, minpay).Scan(&id); err != nil {
		return 0, fmt.Errorf("queryRow: %w", err)
	}
	return id, nil
}

func (r *balance) CreateAccountPool(ctx context.Context, accountID int) error {
	query := `INSERT INTO emcd.accounts_pool (emcd_address_autopay, account_id, created_at, updated_at, autopay_percent)
	SELECT TRUE, $1, NOW(), NOW(), 100
	WHERE NOT EXISTS(
			SELECT NULL
			FROM emcd.accounts_pool
			WHERE (account_id) = ($1)
		)`
	if _, err := r.pool.Exec(ctx, query, accountID); err != nil {
		return fmt.Errorf("queryRow: %w", err)
	}
	return nil
}

func (r *balance) CreateAccountReferral(ctx context.Context, accountID, coinID int) error {
	query := `INSERT INTO emcd.accounts_referral (account_id, coin_id, tier, referral_fee, active_referrals)
	SELECT $1, $2, 1, 0.0005, 0
	WHERE NOT EXISTS(
			SELECT NULL
			FROM emcd.accounts_referral
			WHERE (account_id) = ($1)
		)`
	if _, err := r.pool.Exec(ctx, query, accountID, coinID); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r *balance) FindOperationsAndTransactions(ctx context.Context, queryParams *model.OperationWithTransactionQuery) ([]*model.OperationWithTransaction, int64, error) {

	orderDirection := "DESC"
	if queryParams.Asc {
		orderDirection = "ASC"
	}

	pDateFrom := sql.NullTime{}
	pDateTo := sql.NullTime{}
	pDateTo.Time = time.Now()

	var err error

	if queryParams.DateFrom != "" {
		pDateFrom.Time, err = time.Parse(time.DateTime, queryParams.DateFrom)
		if err != nil {
			return nil, 0, fmt.Errorf("can't parse dateFrom: %w", err)
		}
	}

	if queryParams.DateTo != "" {
		pDateTo.Time, err = time.Parse(time.DateTime, queryParams.DateTo)
		if err != nil {
			return nil, 0, fmt.Errorf("can't parse dateTo: %w", err)
		}
	}

	query := `
		SELECT 
		    o.id,
		    o.account_id,
		    o.coin_id, 
		    t.token_id,
			o.amount,
			o.type,
			o.transaction_id,
			t.action_id,
			t.comment,
			t.fee,
			t.from_referral_id,
			t.gas_price,
			t.hash,
			t.hashrate,
			t.receiver_account_id,
			t.receiver_address,
			t.sender_account_id,
			tb.id AS transaction_block_id,
			tb.blocked_till,
			tb.unblock_to_account_id,
			tb.unblock_transaction_id,
			o.created_at
		FROM emcd.operations o
			 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
			 LEFT JOIN emcd.transactions t ON o.transaction_id = t.id
			 LEFT JOIN emcd.transactions_blocks tb ON tb.block_transaction_id = t.id
		WHERE
		    (ua.user_id = $3 OR Coalesce($3, 0) = 0) AND
		    (o.coin_id = $4 OR Coalesce($4, 0) = 0) AND
			(o.token_id = $5 OR Coalesce($5, 0) = 0) AND
			(t.action_id = $6 OR $6 IS NULL) AND
			(ua.account_type_id = $7 OR Coalesce($7, 0) = 0)AND
			(o.type = ANY($8) OR Coalesce($8,'{}') = '{}') AND
			(o.created_at >= $9 OR $9 IS NULL) AND
			(o.created_at <= $10 OR $10 IS NULL) AND
			(TRUNC(t.amount, 8) = TRUNC($11, 8) OR Coalesce($11, 0) = 0) AND
			(t.hash = $12 OR Coalesce($12, '') = '') AND
			(t.receiver_account_id = $13 OR Coalesce($13, 0) = 0) AND
			(t.receiver_address = $14 OR Coalesce($14, '') = '') AND
			(t.sender_account_id = $15 OR Coalesce($15, 0) = 0) AND
			(tb.id = $16 OR Coalesce($16, 0) = 0) AND
			(tb.unblock_to_account_id = $17 OR Coalesce($17, 0) = 0) AND
			(tb.unblock_transaction_id = $18 OR Coalesce($18, 0) = 0) AND
			(t.from_referral_id = $19 OR Coalesce($19, 0) = 0)
			ORDER BY o.%s %s OFFSET $1 LIMIT $2
			;`

	totalCountQuery := `
		SELECT 
		    count(*)
		FROM emcd.operations o
			 INNER JOIN emcd.users_accounts ua ON ua.id = o.account_id
			 LEFT JOIN emcd.transactions t ON o.transaction_id = t.id
			 LEFT JOIN emcd.transactions_blocks tb ON tb.block_transaction_id = t.id
		WHERE
		    (ua.user_id = $1 OR Coalesce($1, 0) = 0) AND
		    (o.coin_id = $2 OR Coalesce($2, 0) = 0) AND
			(o.token_id = $3 OR Coalesce($3, 0) = 0) AND
			(t.action_id = $4 OR $4 IS NULL) AND
			(ua.account_type_id = $5 OR Coalesce($5, 0) = 0)AND
			(o.type = ANY($6) OR Coalesce($6,'{}') = '{}') AND
			(o.created_at >= $7 OR $7 IS NULL) AND
			(o.created_at <= $8 OR $8 IS NULL) AND
			(TRUNC(t.amount, 8) = TRUNC($9, 8) OR Coalesce($9, 0) = 0) AND
			(t.hash = $10 OR Coalesce($10, '') = '') AND
		    (t.receiver_account_id = $11 OR Coalesce($11, 0) = 0) AND
			(t.receiver_address = $12 OR Coalesce($12, '') = '') AND
			(t.sender_account_id = $13 OR Coalesce($13, 0) = 0) AND
			(tb.id = $14 OR Coalesce($14, 0) = 0) AND
			(tb.unblock_to_account_id = $15 OR Coalesce($15, 0) = 0) AND
			(tb.unblock_transaction_id = $16 OR Coalesce($16, 0) = 0) AND
			(t.from_referral_id = $17 OR Coalesce($17, 0) = 0)
			;`

	result := make([]*model.OperationWithTransaction, 0)

	rows, err := r.pool.Query(ctx, fmt.Sprintf(query, queryParams.SortField, orderDirection), queryParams.Offset, queryParams.Limit,
		queryParams.UserID, queryParams.CoinID, queryParams.TokenID, queryParams.ActionID, queryParams.AccountType, queryParams.OperationTypes,
		pDateFrom.Time, pDateTo.Time, queryParams.Amount, queryParams.Hash, queryParams.ReceiverAccountID, queryParams.ReceiverAddress,
		queryParams.SenderAccountID, queryParams.TransactionBlockID,
		queryParams.UnblockToAccountId, queryParams.UnblockTransactionId, queryParams.FromReferralId)

	if err != nil {
		return nil, 0, fmt.Errorf("can't select: %w", err)
	}

	var count int64
	err = r.pool.QueryRow(ctx, totalCountQuery, queryParams.UserID, queryParams.CoinID, queryParams.TokenID,
		queryParams.ActionID, queryParams.AccountType, queryParams.OperationTypes,
		pDateFrom.Time, pDateTo.Time, queryParams.Amount, queryParams.Hash, queryParams.ReceiverAccountID, queryParams.ReceiverAddress,
		queryParams.SenderAccountID, queryParams.TransactionBlockID,
		queryParams.UnblockToAccountId, queryParams.UnblockTransactionId, queryParams.FromReferralId).Scan(&count)

	if err != nil {
		return nil, 0, fmt.Errorf("can't count select: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			o                    model.OperationWithTransaction
			unblockTransactionId sql.NullInt64
			unblockToAccountID   sql.NullInt64
			fromReferralId       sql.NullInt64
			tokenID              sql.NullInt64
			actionID             sql.NullString
			fee                  decimal.NullDecimal
			gasPrice             decimal.NullDecimal
			hash                 sql.NullString
			receiverAddress      sql.NullString
			hashrate             sql.NullInt64
			transactionBlockID   sql.NullInt64
			blockedTill          sql.NullTime
			comment              sql.NullString
			amount               decimal.NullDecimal
		)

		err = rows.Scan(
			&o.Id,
			&o.AccountID,
			&o.CoinID,
			&tokenID,
			&amount,
			&o.Type,
			&o.TransactionID,
			&actionID,
			&comment,
			&fee,
			&fromReferralId,
			&gasPrice,
			&hash,
			&hashrate,
			&o.ReceiverAccountID,
			&receiverAddress,
			&o.SenderAccountID,
			&transactionBlockID,
			&blockedTill,
			&unblockToAccountID,
			&unblockTransactionId,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("can't scan: %w", err)
		}
		o.UnblockTransactionId = unblockTransactionId.Int64
		o.UnblockToAccountId = unblockToAccountID.Int64
		o.TokenID = tokenID.Int64
		o.ActionID = actionID.String
		o.ReceiverAddress = receiverAddress.String
		o.Hash = hash.String
		o.Fee = fee.Decimal
		o.Hashrate = hashrate.Int64
		o.GasPrice = gasPrice.Decimal
		o.TransactionBlockID = transactionBlockID.Int64
		o.BlockedTill = blockedTill.Time
		o.Comment = comment.String
		o.Amount = amount.Decimal
		result = append(result, &o)
	}

	return result, count, nil
}

func (r *balance) GetBalanceBeforeTransaction(ctx context.Context, accountID, transactionID int64) (decimal.Decimal, error) {

	query := `SELECT COALESCE(SUM(amount), 0) as sum FROM emcd.operations WHERE account_id = $1 AND transaction_id < $2`

	var result decimal.NullDecimal

	row := r.pool.QueryRow(ctx, query,
		accountID,
		transactionID,
	)

	err := row.Scan(&result)

	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("scan GetBalanceBeforeTransaction: %w", err)
	}

	if !result.Valid {
		return decimal.Decimal{}, fmt.Errorf("scan GetBalanceBeforeTransaction validate: %w", pgx.ErrNoRows)
	}

	return result.Decimal, nil
}

func (r *balance) GetP2PAdminId(ctx context.Context) (int, error) {
	q := `select u.id from users u where u.username = 'p2padmin'`

	row := r.pool.QueryRow(ctx, q)

	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}

		return -1, fmt.Errorf("getP2PAdminAccount: %w", err)
	}

	if id == 0 {
		return 0, ErrNotFound
	}

	return id, nil
}
