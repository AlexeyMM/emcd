package repository

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

type PayoutRepository interface {
	GetPayoutsForBlock(ctx context.Context, coinID int64, minPay float32, timestamp time.Time) ([]*model.PayoutForBlock, error)
	GetCurrentPayoutsBlock(ctx context.Context, coinID int64, username string, isService bool) ([]*model.PayoutBlockTransaction, error)
	GetFreePayouts(ctx context.Context, coinID int64) ([]*model.FreePayout, error)
	GetCurrentPayoutsList(ctx context.Context, coinID, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error)
	GetCurrentReferralsPayoutsList(ctx context.Context, coinID, referralId, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error)
	CheckFreePayoutTransaction(ctx context.Context, accountId, transactionId int64) (decimal.Decimal, error)
	CheckPayoutBlockStatus(ctx context.Context, transactionIds []int64) ([]*model.PayoutBlockStatus, error)
	CheckIncomeOperations(ctx context.Context, queryParams model.CheckIncomeOperationsQuery) ([]*model.IncomeWithFee, error)
	CheckOthers(ctx context.Context, queryParams model.CheckOtherQuery) ([]*model.OtherOperationsWithTransaction, error)
	GetAveragePaid(ctx context.Context, queryParams model.AveragePaidQuery) (decimal.Decimal, error)
	GetServiceUserData(ctx context.Context, coinId int64, username string, limit int64) ([]*model.ServiceUserBlock, error)
}

type payouts struct {
	pool *pgxpool.Pool
}

func NewPayouts(pool *pgxpool.Pool) PayoutRepository {
	return &payouts{
		pool: pool,
	}
}

func (r *payouts) GetPayoutsForBlock(ctx context.Context, coinID int64, minPay float32, timestamp time.Time) ([]*model.PayoutForBlock, error) {
	query := `SELECT acc_id as account_id, userid as user_id, acc_balance as balance, pay_addresses as addresses FROM emcd.getPayoutsListNew($1,$2,$3)`
	rows, err := r.pool.Query(ctx, query, coinID, timestamp.Unix(), minPay)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.PayoutForBlock, 0)

	for rows.Next() {
		var (
			p          model.PayoutForBlock
			account_id sql.NullInt64
			user_id    sql.NullInt64
			balance    decimal.NullDecimal
			addresses  sql.NullString
		)

		err = rows.Scan(
			&account_id,
			&user_id,
			&balance,
			&addresses,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.AccountID = account_id.Int64
		p.UserID = user_id.Int64
		p.Balance = balance.Decimal
		p.Address = addresses.String

		result = append(result, &p)
	}

	return result, nil
}

func (r *payouts) GetCurrentPayoutsBlock(ctx context.Context, coinID int64, username string, isService bool) ([]*model.PayoutBlockTransaction, error) {

	query := `SELECT t.id, t.amount as balance FROM emcd.transactions t
  LEFT JOIN emcd.transactions_blocks tb ON t.id = tb.block_transaction_id
  JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
  JOIN emcd.users u ON ua.user_id = u.id
  LEFT JOIN emcd.service_users su ON u.id = su.user_id
  WHERE t.type = $1 AND t.coin_id = $2 AND tb.id IS NOT NULL AND tb.unblock_transaction_id IS NULL AND u.nopay = false`

	if username == "" {
		query = query + ` AND u.username = '` + username + `'`
	}

	if isService {
		query = fmt.Sprintf("%s AND su.user_id IS NOT NULL", query)
	} else {
		query = fmt.Sprintf("%s AND su.user_id IS NULL", query)
	}

	rows, err := r.pool.Query(ctx, query, model.PoolPaysUsersBalanceTrTypeID, coinID)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.PayoutBlockTransaction, 0)

	for rows.Next() {
		var (
			p       model.PayoutBlockTransaction
			id      sql.NullInt64
			balance decimal.NullDecimal
		)

		err = rows.Scan(
			&id,
			&balance,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.ID = id.Int64
		p.Balance = balance.Decimal

		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) GetFreePayouts(ctx context.Context, coinID int64) ([]*model.FreePayout, error) {

	query := `SELECT 
    t.amount,
    t.action_id,
    t.coin_id,
    t.comment,
    t.created_at,
    t.fee,
    t.from_referral_id,
    t.gas_price,
    t.hash,
    t.hashrate,
    t.id,
    t.is_viewed,
    t.receiver_account_id,
    t.receiver_address,
    t.sender_account_id,
    t.token_id,
    t.type,
    ua.id as account_id, 
    u.id as user_id, 
    u.username as username 
FROM emcd.transactions_blocks bt
LEFT JOIN emcd.transactions t ON bt.block_transaction_id = t.id
LEFT JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
LEFT JOIN emcd.users u ON ua.user_id = u.id
WHERE bt.unblock_transaction_id is null and t.receiver_address is not null and t.type = $1 and u.nopay = false
and t.coin_id = $2`

	rows, err := r.pool.Query(ctx, query, model.WalletMiningTransferTrTypeID, coinID)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.FreePayout, 0)

	for rows.Next() {
		var (
			p                 model.FreePayout
			userId            sql.NullInt64
			username          sql.NullString
			amount            decimal.NullDecimal
			comment           sql.NullString
			createdAt         sql.NullTime
			fee               decimal.NullDecimal
			fromReferralId    sql.NullInt64
			gasPrice          decimal.NullDecimal
			hash              sql.NullString
			hashrate          sql.NullInt64
			isViewed          sql.NullBool
			receiverAccountID sql.NullInt64
			receiverAddress   sql.NullString
			senderAccountID   sql.NullInt64
			tokenID           sql.NullInt64
			coinId            sql.NullInt64
			actionId          sql.NullString
		)

		err = rows.Scan(
			&amount,
			&actionId,
			&coinId,
			&comment,
			&createdAt,
			&fee,
			&fromReferralId,
			&gasPrice,
			&hash,
			&hashrate,
			&p.ID,
			&isViewed,
			&receiverAccountID,
			&receiverAddress,
			&senderAccountID,
			&tokenID,
			&p.Type,
			&p.AccountId,
			&userId,
			&username,
		)

		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.Amount = amount.Decimal
		p.UserId = userId.Int64
		p.Username = username.String
		p.Comment = comment.String
		p.CreatedAt = createdAt.Time
		p.Fee = fee.Decimal
		p.FromReferralID = fromReferralId.Int64
		p.GasPrice = gasPrice.Decimal
		p.Hash = hash.String
		p.Hashrate = hashrate.Int64
		p.IsViewer = isViewed.Bool
		p.ReceiverAccountID = receiverAccountID.Int64
		p.ReceiverAddress = receiverAddress.String
		p.SenderAccountID = senderAccountID.Int64
		p.TokenID = tokenID.Int64
		p.CoinID = coinId.Int64
		p.ActionID = actionId.String
		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) GetCurrentPayoutsList(ctx context.Context, coinID, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error) {

	query := `SELECT 
    ua.id AS id, 
    ua.id AS account_id2, 
    u.id AS user_id, 
    u.ref_id AS ref_id, 
    ua.coin_id, 
    u.username, 
    ua.minpay, 
    u.master_id,  
    block.address AS address,
    block.amount AS balance, 
    block.block_tr_id AS block_id, 
    block.created_at AS block_create,
    /** Запрос конрольных значений */
       ( SELECT json_build_object(
                    /** Получаем вывод с коинхолд */
                        'coinhold', coalesce(SUM(CASE WHEN op_m.type IN ($1) THEN op_m.amount else 0 end), 0),
                    /** Получаем начисления за майнинг */
                        'incomes', coalesce(SUM(CASE WHEN op_m.type IN ($2) THEN op_m.amount else 0 end), 0),
                    /** Получаем хэшрейт */
                        'hashrate' , coalesce(SUM(CASE WHEN op_m.type = $2 THEN t.hashrate else 0 end), 0),
                    /** Получаем списания со счета. Комиссия, донаты, увод на депозит */
                        'feeAndMore', coalesce(SUM(CASE WHEN op_m.type IN ($3, $4, $5) THEN op_m.amount else 0 end), 0),
                    /** Получаем реферальное вознаграждение */
                        'ref', coalesce(SUM(CASE WHEN op_m.type IN ($6) THEN op_m.amount else 0 end), 0),
                    /** Получаем сумму операций, которые не попали в предыдущие выборки, кроме 21 блокировки */
                        'other', coalesce(SUM(CASE WHEN op_m.type NOT IN ($6, $1, $5, $2, $3, $4, $7) THEN op_m.amount else 0 end),0),
                    /** Получаем типы транзакций которые попали в other, кроме 21 блокировки */
                        'types', string_agg((CASE WHEN op_m.type NOT IN ($6, $1, $5, $2, $3, $4, $7) THEN op_m.type end)::character varying,',' ORDER BY op_m.type),
                        'accountID', op_m.account_id,
                     'lastPay', coalesce(pay.last_pay, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE,
                     'incomeFirst',  coalesce(income.first, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE,
                     'incomeLast', coalesce(income.last, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE
                    )
         	FROM emcd.operations op_m
                  LEFT JOIN emcd.users_accounts ua_d ON op_m.account_id = ua_d.id 
                  LEFT JOIN (SELECT (max(o.created_at)) as last_pay, ua.id as account_id  FROM emcd.transactions tb
                  LEFT JOIN emcd.transactions_blocks b ON b.block_transaction_id = tb.id
                  LEFT JOIN emcd.transactions tu ON b.unblock_transaction_id = tu.id
                  LEFT JOIN emcd.users_accounts ua ON tb.sender_account_id = ua.id
                  LEFT JOIN emcd.operations o ON o.transaction_id = tb.id   
                  WHERE ((tu.id is not null and tu.receiver_account_id != tb.sender_account_id) or b.id IS NULL) and tb.type = $8 and tb.coin_id = $9
                  GROUP BY ua.id) pay ON pay.account_id = ua_d.id
                  LEFT JOIN (
                  	SELECT 
                  	    to_timestamp(to_char(min(created_at - INTERVAL '1 DAY'), 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM-DD HH24:MI:SS') as first, 
                  	    to_timestamp(to_char(max(created_at), 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM-DD HH24:MI:SS') as last,
                  	    account_id 
                  	FROM emcd.operations WHERE type = $2 GROUP BY account_id) income ON income.account_id = ua_d.id
                  LEFT JOIN emcd.transactions t ON t.id = op_m.transaction_id
         WHERE op_m.created_at > coalesce(pay.last_pay, '1970-01-01 00:00:00') and op_m.created_at < block.created_at and op_m.account_id = ua.id
         GROUP BY op_m.account_id, last_pay, income.first, income.last
       ) as calc
FROM emcd.users_accounts ua
         LEFT JOIN emcd.users u on u.id = ua.user_id
         LEFT JOIN emcd.service_users su ON u.id = su.user_id --исключаем попадание сервисных юзеров
         LEFT JOIN (SELECT  t.id as block_tr_id, ua.id as account_id, t.amount, t.created_at as created_at, t.receiver_address as address FROM emcd.transactions_blocks bt
         LEFT JOIN emcd.transactions t ON bt.block_transaction_id = t.id
         LEFT JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
         WHERE bt.unblock_transaction_id is null and t.receiver_address is not null and t.type = $8 and t.coin_id = $9) block ON block.account_id = ua.id
WHERE block.account_id IS NOT NULL and su.user_id IS NULL;`

	rows, err := r.pool.Query(ctx, query,
		model.CnhldCloseTrTypeID,
		model.MainCoinMiningPayoutTrTypeID,
		model.UserPaysPoolComsTrTypeID,
		model.UserPaysPoolDonationsTrTypeID,
		model.FromMiningReplenishmentTrTypeID,
		model.PayUserReferralBonusTrTypeID,
		model.PoolPaysUsersBalanceTrTypeID,
		paymentTransactionType,
		coinID)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.PayoutWithCalculation, 0)

	type TmpCalc struct {
		Coinhold    decimal.NullDecimal `json:"coinhold"`
		Incomes     decimal.NullDecimal `json:"incomes"`
		Hashrate    decimal.NullDecimal `json:"hashrate"`
		FeeAndMore  decimal.NullDecimal `json:"feeAndMore"`
		Ref         decimal.NullDecimal `json:"ref"`
		Other       decimal.NullDecimal `json:"other"`
		Types       string              `json:"types"`
		AccountID   int64               `json:"accountID"`
		LastPay     time.Time           `json:"lastPay"`
		IncomeFirst time.Time           `json:"incomeFirst"`
		IncomeLast  time.Time           `json:"incomeLast"`
	}

	for rows.Next() {
		var (
			p              model.PayoutWithCalculation
			refId          sql.NullInt64
			minpay         decimal.NullDecimal
			amount         decimal.NullDecimal
			address        sql.NullString
			blockId        sql.NullInt64
			masterId       sql.NullInt64
			blockCreatedAt sql.NullTime
			tmp            []byte
		)

		err = rows.Scan(
			&p.ID,
			&p.AccountID2,
			&p.UserID,
			&refId,
			&p.CoinID,
			&p.Username,
			&minpay,
			&masterId,
			&address,
			&amount,
			&blockId,
			&blockCreatedAt,
			&tmp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		var calc TmpCalc
		if err = json.Unmarshal(tmp, &calc); err != nil {
			return nil, fmt.Errorf("calc, unmarshal: %w", err)
		}

		p.Minpay = minpay.Decimal
		p.RefID = refId.Int64
		p.Address = address.String
		p.Balance = amount.Decimal
		p.BlockID = blockId.Int64
		p.BlockCreate = blockCreatedAt.Time
		p.MasterID = masterId.Int64

		p.Calc.Coinhold = calc.Coinhold.Decimal
		p.Calc.Incomes = calc.Incomes.Decimal
		p.Calc.Hashrate = calc.Hashrate.Decimal
		p.Calc.FeeAndMore = calc.FeeAndMore.Decimal
		p.Calc.Ref = calc.Ref.Decimal
		p.Calc.Other = calc.Other.Decimal
		p.Calc.Types = calc.Types
		p.Calc.AccountID = calc.AccountID
		p.Calc.LastPay = calc.LastPay
		p.Calc.IncomeFirst = calc.IncomeFirst
		p.Calc.IncomeLast = calc.IncomeLast

		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) GetCurrentReferralsPayoutsList(ctx context.Context, coinID, referralId, paymentTransactionType int64) ([]*model.PayoutWithCalculation, error) {

	query := `SELECT 
    ua.id AS id, 
    ua.id as account_id2, 
    u.id as user_id, 
    u.ref_id as ref_id, 
    ua.coin_id, 
    u.username, 
    ua.minpay, 
    u.master_id,  
    ua.address as address, 
    (SELECT sum(amount) FROM emcd.operations o where o.account_id = ua.id) AS balance,
    /** Запрос конрольных значений */
       ( SELECT json_build_object(
                    /** Получаем вывод с коинхолд */
                        'coinhold', coalesce(SUM(CASE WHEN op_m.type IN ($1) THEN op_m.amount else 0 end), 0),
                    /** Получаем начисления за майнинг */
                    'incomes', coalesce(SUM(CASE WHEN op_m.type IN ($2) THEN op_m.amount else 0 end), 0),
                    /** Получаем хэшрейт */
                    'hashrate' , coalesce(SUM(CASE WHEN op_m.type = $2 THEN t.hashrate else 0 end), 0),
                    /** Получаем списания со счета. Комиссия, донаты, увод на депозит */
                    'feeAndMore', coalesce(SUM(CASE WHEN op_m.type IN ($3, $4, $5) THEN op_m.amount else 0 end), 0),
                    /** Получаем реферальное вознаграждение */
                    'ref', coalesce(SUM(CASE WHEN op_m.type IN ($6) THEN op_m.amount else 0 end), 0),
                    /** Получаем сумму операций, которые не попали в предыдущие выборки, кроме 21 блокировки */
                    'other', coalesce(SUM(CASE WHEN op_m.type NOT IN ($6, $1, $5, $2, $3, $4, $7) THEN op_m.amount else 0 end),0),
                    /** Получаем типы транзакций которые попали в other, кроме 21 блокировки */
                    'types', string_agg((CASE WHEN op_m.type NOT IN ($6, $1, $5, $2, $3, $4, $7) THEN op_m.type end)::character varying,',' ORDER BY op_m.type),
                    'accountID', op_m.account_id,
                    'lastPay', coalesce(pay.last_pay, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE,
                    'incomeFirst', coalesce(income.first, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE,
                    'incomeLast', coalesce(income.last, '1970-01-01 00:00:00')::TIMESTAMP WITH TIME ZONE
                    )
         FROM emcd.operations op_m
                  LEFT JOIN emcd.users_accounts ua_d ON op_m.account_id = ua_d.id
                  LEFT JOIN (SELECT max(created_at) as last_pay, account_id FROM emcd.operations WHERE type = $8 GROUP BY account_id) pay ON pay.account_id = ua_d.id
                  LEFT JOIN (SELECT min(created_at - INTERVAL '1 DAY') as first, max(created_at) as last, account_id 
                             FROM emcd.operations WHERE type = $2 GROUP BY account_id) income ON income.account_id = ua_d.id
                  LEFT JOIN emcd.transactions t ON t.id = op_m.transaction_id
         WHERE op_m.created_at > coalesce(pay.last_pay, '1970-01-01 00:00:00') and op_m.account_id = ua.id
         GROUP BY op_m.account_id, last_pay, income.first, income.last
       ) as calc
FROM emcd.users_accounts ua
         LEFT JOIN emcd.users u on u.id = ua.user_id
         LEFT JOIN emcd.service_users su ON u.id = su.user_id
where ua.account_type_id = 2
  and u.ref_id = $9
  and (SELECT sum(amount) FROM emcd.operations o where o.account_id = ua.id) >= 0.00000001
  and coin_id = $10
  and nopay = false
  and u.is_autopay_disabled = false
  and su.user_id IS NULL`

	rows, err := r.pool.Query(ctx, query,
		model.CnhldCloseTrTypeID,
		model.MainCoinMiningPayoutTrTypeID,
		model.UserPaysPoolComsTrTypeID,
		model.UserPaysPoolDonationsTrTypeID,
		model.FromMiningReplenishmentTrTypeID,
		model.PayUserReferralBonusTrTypeID,
		model.PoolPaysUsersBalanceTrTypeID,
		paymentTransactionType,
		referralId,
		coinID)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.PayoutWithCalculation, 0)

	type TmpCalc struct {
		Coinhold    decimal.NullDecimal `json:"coinhold"`
		Incomes     decimal.NullDecimal `json:"incomes"`
		Hashrate    decimal.NullDecimal `json:"hashrate"`
		FeeAndMore  decimal.NullDecimal `json:"feeAndMore"`
		Ref         decimal.NullDecimal `json:"ref"`
		Other       decimal.NullDecimal `json:"other"`
		Types       string              `json:"types"`
		AccountID   int64               `json:"accountID"`
		LastPay     time.Time           `json:"lastPay"`
		IncomeFirst time.Time           `json:"incomeFirst"`
		IncomeLast  time.Time           `json:"incomeLast"`
	}

	for rows.Next() {
		var (
			p         model.PayoutWithCalculation
			refId     sql.NullInt64
			minpay    decimal.NullDecimal
			amount    decimal.NullDecimal
			address   sql.NullString
			tmp       []byte
			master_id sql.NullInt64
		)

		err = rows.Scan(
			&p.ID,
			&p.AccountID2,
			&p.UserID,
			&refId,
			&p.CoinID,
			&p.Username,
			&minpay,
			&master_id,
			&address,
			&amount,
			&tmp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		var calc TmpCalc
		if err = json.Unmarshal(tmp, &calc); err != nil {
			return nil, fmt.Errorf("calc, unmarshal: %w", err)
		}

		p.Minpay = minpay.Decimal
		p.RefID = refId.Int64
		p.Address = address.String
		p.Balance = amount.Decimal
		p.MasterID = master_id.Int64

		p.Calc.Coinhold = calc.Coinhold.Decimal
		p.Calc.Incomes = calc.Incomes.Decimal
		p.Calc.Hashrate = calc.Hashrate.Decimal
		p.Calc.FeeAndMore = calc.FeeAndMore.Decimal
		p.Calc.Ref = calc.Ref.Decimal
		p.Calc.Other = calc.Other.Decimal
		p.Calc.Types = calc.Types
		p.Calc.AccountID = calc.AccountID
		p.Calc.LastPay = calc.LastPay
		p.Calc.IncomeFirst = calc.IncomeFirst
		p.Calc.IncomeLast = calc.IncomeLast

		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) CheckFreePayoutTransaction(ctx context.Context, accountId, transactionId int64) (decimal.Decimal, error) {
	query := `SELECT COALESCE(SUM(amount), 0) as sum FROM emcd.operations WHERE account_id = $1 AND transaction_id < $2`

	var result decimal.NullDecimal
	row := r.pool.QueryRow(ctx, query, accountId, transactionId)

	err := row.Scan(&result)

	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("scan CheckFreePayoutTransaction: %w", err)
	}
	if !result.Valid {
		return decimal.Decimal{}, fmt.Errorf("scan CheckFreePayoutTransaction: %w", pgx.ErrNoRows)
	}

	return result.Decimal, nil
}

func (r *payouts) CheckPayoutBlockStatus(ctx context.Context, transactionIds []int64) ([]*model.PayoutBlockStatus, error) {

	transactionIdsString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(transactionIds)), ","), "[]")

	query := fmt.Sprintf("SELECT bt.unblock_to_account_id as to_account_id, "+
		"t.type, t.receiver_address, "+
		"bt.unblock_transaction_id as ub_tr_id, "+
		"t.amount "+
		"FROM emcd.transactions_blocks bt "+
		"LEFT JOIN emcd.transactions t ON bt.block_transaction_id = t.id "+
		"WHERE bt.block_transaction_id IN (%s)", transactionIdsString)

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.PayoutBlockStatus, 0)

	for rows.Next() {
		var (
			p                    model.PayoutBlockStatus
			toAccountId          sql.NullInt64
			trType               sql.NullInt64
			receiverAddress      sql.NullString
			unblockTransactionId sql.NullInt64
			amount               decimal.NullDecimal
		)

		err = rows.Scan(
			&toAccountId,
			&trType,
			&receiverAddress,
			&unblockTransactionId,
			&amount,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.ToAccountId = toAccountId.Int64
		p.Type = trType.Int64
		p.ReceiverAddress = receiverAddress.String
		p.UnblockTransactionId = unblockTransactionId.Int64
		p.Amount = amount.Decimal

		result = append(result, &p)
	}

	return result, nil
}

func (r *payouts) CheckIncomeOperations(ctx context.Context, queryParams model.CheckIncomeOperationsQuery) ([]*model.IncomeWithFee, error) {
	query := `
		SELECT t.id as transaction_id, o.amount as amount,
			   t.hashrate as hashrate,
			   o.created_at as created_at,
			   (case when o.created_at < $1 then (SELECT fee FROM histories.userfee_history WHERE coin = $2
											and userid = $3
											and date_start <= to_date(cast(o.created_at as varchar), '%Y-%MM-%DD')
											and (date_end >= to_date(cast(o.created_at as varchar), '%Y-%MM-%DD') OR date_end is NULL) LIMIT 1
			   ) else null end) as fee
			   
		from emcd.operations o
		LEFT JOIN emcd.transactions t ON t.id = o.transaction_id
		WHERE account_id = $4 and o.type = $5 and o.created_at >= $6
		ORDER BY o.created_at ASC`

	rows, err := r.pool.Query(ctx, query,
		queryParams.CreatedAt.String(),
		queryParams.Coin,
		queryParams.UserID,
		queryParams.AccountID,
		model.MainCoinMiningPayoutTrTypeID,
		queryParams.LastPayAt.String())

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.IncomeWithFee, 0)

	for rows.Next() {
		var (
			p             model.IncomeWithFee
			transactionId sql.NullInt64
			amount        decimal.NullDecimal
			hashrate      sql.NullInt64
			createdAt     sql.NullTime
			fee           decimal.NullDecimal
		)

		err = rows.Scan(
			&transactionId,
			&amount,
			&hashrate,
			&createdAt,
			fee,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.TransactionId = transactionId.Int64
		p.Amount = amount.Decimal
		p.Hashrate = hashrate.Int64
		p.CreatedAt = createdAt.Time
		p.Fee = fee.Decimal

		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) GetAveragePaid(ctx context.Context, queryParams model.AveragePaidQuery) (decimal.Decimal, error) {

	interval := strconv.FormatInt(queryParams.Days, 10) + " day"

	query := fmt.Sprintf(`SELECT COALESCE(AVG(sum), 0) AS avg FROM (
		SELECT SUM(amount) AS sum, DATE(created_at) AS date FROM emcd.transactions WHERE
		coin_id = $1 AND type = $2 and created_at  >= (current_date - INTERVAL '%s')::date and date(created_at) < date(NOW())
        AND receiver_account_id = (
		  SELECT ua.id FROM emcd.users_accounts ua
		  JOIN emcd.users u ON ua.user_id = u.id
		  WHERE u.username = $3 AND account_type_id = $4 AND coin_id = $1 LIMIT 1)
		  GROUP BY DATE(created_at)
		) as t`, interval)

	var result decimal.NullDecimal

	row := r.pool.QueryRow(ctx, query,
		queryParams.CoinID,
		queryParams.TransactionTypeID,
		queryParams.Username,
		queryParams.AccountTypeID,
	)

	err := row.Scan(&result)

	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("scan GetAveragePaid: %w", err)
	}

	if !result.Valid {
		return decimal.Decimal{}, fmt.Errorf("scan GetAveragePaid validate: %w", pgx.ErrNoRows)
	}

	return result.Decimal, nil
}

func (r *payouts) CheckOthers(ctx context.Context, queryParams model.CheckOtherQuery) ([]*model.OtherOperationsWithTransaction, error) {

	query := `SELECT t.id AS transaction_id,
       t.sender_account_id AS sender_id,
       t.receiver_account_id AS receiver_id,
       t.hash AS hash,
       o.id AS operation_id,
       ua_s.user_id AS sender_user_id,
       ua_r.user_id AS receiver_user_id,
       o.amount AS amount,
       o.type AS type,
       o.created_at AS created_at,
       t.comment AS comment
FROM emcd.operations o
LEFT JOIN emcd.transactions t ON t.id = o.transaction_id
LEFT JOIN emcd.users_accounts ua_s ON ua_s.id = t.sender_account_id
LEFT JOIN emcd.users_accounts ua_r ON ua_r.id = t.receiver_account_id
WHERE 
    o.account_id = $1 AND 
    o.type = ANY ($2) AND
    o.created_at >= $3 AND
    (o.created_at < $4 OR $4 IS NULL)
ORDER BY o.created_at ASC`

	rows, err := r.pool.Query(ctx, query, queryParams.AccountID, queryParams.Types, queryParams.LastPayAt, queryParams.BlockCreatedAt)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.OtherOperationsWithTransaction, 0)

	for rows.Next() {
		var (
			p               model.OtherOperationsWithTransaction
			transactionId   sql.NullInt64
			senderId        sql.NullInt64
			receiverId      sql.NullInt64
			hash            sql.NullString
			operationId     sql.NullInt64
			senderUserId    sql.NullInt64
			receiverUserId  sql.NullInt64
			amount          decimal.NullDecimal
			transactionType sql.NullInt64
			createdAt       sql.NullTime
			comment         sql.NullString
		)

		err = rows.Scan(
			&transactionId,
			&senderId,
			&receiverId,
			&hash,
			&operationId,
			&senderUserId,
			&receiverUserId,
			&amount,
			&transactionType,
			&createdAt,
			&comment,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.TransactionID = transactionId.Int64
		p.SenderID = senderId.Int64
		p.ReceiverID = receiverId.Int64
		p.Hash = hash.String
		p.OperationID = operationId.Int64
		p.SenderUserID = senderUserId.Int64
		p.ReceiverUserID = receiverUserId.Int64
		p.Amount = amount.Decimal
		p.TransactionTypeID = transactionType.Int64
		p.CreatedAt = createdAt.Time
		p.Comment = comment.String

		result = append(result, &p)
	}

	return result, nil

}

func (r *payouts) GetServiceUserData(ctx context.Context, coinId int64, username string, limit int64) ([]*model.ServiceUserBlock, error) {

	query := `SELECT 
    t.receiver_address as address, 
    ua.id as su_account_id, 
    u.id as user_id, 
    u.username as username, 
    t.amount as amount,
    t.id as block_id 
FROM emcd.transactions_blocks bt
LEFT JOIN emcd.transactions t ON bt.block_transaction_id = t.id
LEFT JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
LEFT JOIN emcd.users u ON ua.user_id = u.id
LEFT JOIN emcd.service_users su ON u.id = su.user_id
WHERE bt.unblock_transaction_id is null and t.receiver_address is not null and t.type = $1 and u.nopay = false and u.username = $2 AND su.user_id is not null
and t.coin_id = $3 ORDER BY t.created_at DESC LIMIT $4`

	rows, err := r.pool.Query(ctx, query, model.PoolPaysUsersBalanceTrTypeID, username, coinId, limit)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.ServiceUserBlock, 0)

	for rows.Next() {
		var (
			p           model.ServiceUserBlock
			address     sql.NullString
			suAccountId sql.NullInt64
			userId      sql.NullInt64
			uname       sql.NullString
			amount      decimal.NullDecimal
			blockId     sql.NullInt64
		)

		err = rows.Scan(
			&address,
			&suAccountId,
			&userId,
			&uname,
			&amount,
			&blockId,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.Address = address.String
		p.SuAccountID = suAccountId.Int64
		p.UserID = userId.Int64
		p.Username = uname.String
		p.Amount = amount.Decimal
		p.BlockID = blockId.Int64

		result = append(result, &p)
	}

	return result, nil
}
