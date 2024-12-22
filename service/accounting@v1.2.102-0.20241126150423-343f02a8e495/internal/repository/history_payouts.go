package repository

import (
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"context"
	"fmt"
	"github.com/shopspring/decimal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"code.emcdtech.com/emcd/service/accounting/model"
)

const (
	selectPayoutsNew = `
	SELECT
		CAST(EXTRACT(EPOCH FROM t.created_at) AS INT) AS time,
		TRUNC(
			COALESCE(ABS(amount), 0),
			8) AS amount,
		hash AS tx
	FROM
		emcd.transactions t
	LEFT JOIN
		emcd.users_accounts ua ON t.sender_account_id = ua.id
	WHERE
	  ua.user_id = $1
	  AND (ua.account_type_id = 2 OR ua.account_type_id = 1)
	  AND ua.coin_id = $2`

	selectPayoutsIsViewed = `
	SELECT CASE WHEN count(*) > 0 THEN true ELSE false END AS temp
	FROM transactions t
			 LEFT JOIN users_accounts ua ON ua.id = t.sender_account_id
	WHERE t.is_viewed = false
	  AND ua.user_id = $1
	  AND ua.coin_id = $2
	  AND t.type = 21
	  AND ua.account_type_id = 2`
)

type PayoutsHistory interface {
	GetNewPayouts(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error)
	GetPayoutsIsViewed(ctx context.Context, id int, coinCode string) (bool, error)
}

type payoutsHistory struct {
	pool          *pgxpool.Pool
	slaveDBPool   *pgxpool.Pool
	coinValidator coinValidatorRepo.CoinValidatorRepository
}

func NewPayoutsHistory(
	pool *pgxpool.Pool,
	slaveDBPool *pgxpool.Pool,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) PayoutsHistory {
	return &payoutsHistory{
		pool:          pool,
		slaveDBPool:   slaveDBPool,
		coinValidator: coinValidator,
	}
}

func (m *payoutsHistory) GetNewPayouts(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error) {
	coinID, ok := m.coinValidator.GetIdByCode(filter.CoinCode)
	if !ok {
		return nil, fmt.Errorf("failed to find coinId by code: %s", filter.CoinCode)
	}

	var filterValues = []any{filter.UserID, coinID}

	values, query, countSumQuery := getNewPayoutsFilterQuery(filterValues, selectPayoutsNew, filter)
	rows, err := m.slaveDBPool.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("selectPayoutsNew Query: %w: internal error 0002-01", err)
	}

	defer rows.Close()

	payoutsRes := &model.HistoryOutput{
		TotalCount:    0,
		IncomesSum:    nil,
		PayoutsSum:    utils.DecimalToPtr(decimal.NewFromInt(0)),
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       make([]*model.Payout, 0, filter.Limit),
		Wallets:       nil,
	}

	for rows.Next() {
		payout := new(model.Payout)

		if err = rows.Scan(
			&payout.Time,
			&payout.Amount,
			&payout.Tx,
		); err != nil {
			return nil, fmt.Errorf("selectPayoutsNew Scan: %w: internal error 0002-02", err)
		}
		payout.TxId = payout.Tx

		payoutsRes.Payouts = append(payoutsRes.Payouts, payout)
	}

	if len(values) > veryCoolNumber2 {
		values = values[:len(values)-veryCoolNumber2]
	}

	countRows, err := m.slaveDBPool.Query(ctx, countSumQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("countSumQuery Query: %w", err)
	}

	defer countRows.Close()

	if countRows.Next() {
		if err = countRows.Scan(&payoutsRes.TotalCount, &payoutsRes.PayoutsSum); err != nil {
			return nil, fmt.Errorf("countSumQuery Scan: %w", err)
		}
	}

	return payoutsRes, nil
}

func (m *payoutsHistory) GetPayoutsIsViewed(ctx context.Context, id int, coinCode string) (bool, error) {
	if coinCode == "" {
		coinCode = coinBTC
	}
	coinId, ok := m.coinValidator.GetIdByCode(coinCode)
	if !ok {
		return false, fmt.Errorf("failed to find coinId by code: %s", coinCode)
	}
	var isViewed bool
	if err := m.pool.QueryRow(ctx, selectPayoutsIsViewed, id, coinId).
		Scan(&isViewed); err != nil {
		return false, fmt.Errorf("QueryRow: %w", err)
	}
	return isViewed, nil
}

func getNewPayoutsFilterQuery(values []any, query string, filter *model.HistoryInput) ([]any, string, string) {
	if filter.From != "" {
		values = append(values, filter.From)
		query += fmt.Sprintf(` AND date(t.created_at) >= $%d `, len(values))
	}

	if filter.To != "" {
		values = append(values, filter.To)
		query += fmt.Sprintf(` AND date(t.created_at) <= $%d `, len(values))
	}

	values = append(values, pq.Array(filter.TransactionTypesIDs))

	query += fmt.Sprintf(` AND t.type = ANY ($%d)`, len(values))

	query += " \n\tORDER BY t.created_at DESC\n\t"

	countSumQuery := "SELECT COUNT(time), COALESCE(SUM(amount), 0) FROM (" + query + ") count_sum_table"

	values = append(values, filter.Limit)
	query += fmt.Sprintf(` LIMIT $%d`, len(values))

	values = append(values, filter.Offset)
	query += fmt.Sprintf(` OFFSET $%d`, len(values))

	return values, query, countSumQuery
}
