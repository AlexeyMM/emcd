package repository

import (
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"code.emcdtech.com/emcd/service/accounting/internal/controller"
	"code.emcdtech.com/emcd/service/accounting/model"
)

const veryCoolNumber2 = 2

const mainCoinMiningPayoutTrTypeIDTypeCast = 1
const poolPaysUserCompensationTrTypeIDTypeCast = 2
const payUserReferralBonusTrTypeIDTypeCast = 3
const userPaysPoolDonationsTrTypeIDTypeCast = 4
const payoutTrTypeIDTypeCast = 8
const fromMiningReplenishmentTrTypeIDTypeCast = 9
const accrualsChineseTechAccountTrTypeIDTypeCast = 33

const (
	selectComplexityIncomes = `
	SELECT
		date(t.created_at) as date,
		max(extract (epoch from t.created_at)) as time
	FROM
		transactions t
	WHERE
	  	t.receiver_account_id IN (SELECT id FROM users_accounts WHERE user_id = $1 AND coin_id = $2)
		AND t.type = 5
	  	AND t.coin_id = $2
		AND t.amount != 0
	GROUP BY 
		date
	ORDER BY
		time`

	selectIncomesNew = `
WITH user_accounts AS (
    SELECT id
    FROM users_accounts
    WHERE user_id = $1 AND coin_id = $2
),
daily_summary AS (
    SELECT
        date(t.created_at) as date,
        t.type,
        MAX(EXTRACT(EPOCH FROM t.created_at)) as time,
        SUM(COALESCE(t.amount, 0)) as total_amount,
        SUM(COALESCE(t.hashrate, 0)) as total_hashrate
    FROM
        transactions t
    JOIN user_accounts ua ON t.receiver_account_id = ua.id
    WHERE
        t.coin_id = $2
        AND t.amount != 0
        %s
    GROUP BY date, t.type
),
daily_summary_fees AS (
    SELECT
        date(t2.created_at) as date,
        SUM(ABS(COALESCE(t2.amount, 0))) as fee_amount
    FROM
        transactions t2
    JOIN user_accounts ua ON t2.sender_account_id = ua.id
    WHERE
        t2.coin_id = $2
        AND t2.type IN (7, 23)
    GROUP BY date
)
SELECT
    dt.date,
    dt.type,
    dt.time,
    TRUNC(
        dt.total_amount - COALESCE(df.fee_amount, 0), 8
    ) as amount,
    dt.total_hashrate as hashrate
FROM
    daily_summary dt
LEFT JOIN daily_summary_fees df
    ON dt.date = df.date AND dt.type = 5
ORDER BY dt.date DESC
%s`

	selectIncomesIsViewed = `
	SELECT CASE WHEN count(*) > 0 THEN true ELSE false END AS temp
	FROM emcd.transactions t
			 LEFT JOIN emcd.users_accounts ua ON ua.id = t.receiver_account_id
	WHERE t.is_viewed = false
	  AND ua.user_id = $1
	  AND ua.coin_id = $2
	  AND t.type IN (3, 4, 5, 6, 7, 8, 9, 22, 23, 13)
	  AND ua.account_type_id = 2`
)

func incomeTypeCast(transactionType int64) int {
	switch transactionType {
	case model.MainCoinMiningPayoutTrTypeID.Int64():
		return mainCoinMiningPayoutTrTypeIDTypeCast
	case model.PoolPaysUserCompensationTrTypeID.Int64():
		return poolPaysUserCompensationTrTypeIDTypeCast
	case model.PayUserReferralBonusTrTypeID.Int64():
		return payUserReferralBonusTrTypeIDTypeCast
	case model.UserPaysPoolDonationsTrTypeID.Int64():
		return userPaysPoolDonationsTrTypeIDTypeCast
	case model.PayoutTrTypeID.Int64():
		return payoutTrTypeIDTypeCast
	case model.FromMiningReplenishmentTrTypeID.Int64():
		return fromMiningReplenishmentTrTypeIDTypeCast
	case model.AccrualsChineseTechAccountTrTypeID.Int64():
		return accrualsChineseTechAccountTrTypeIDTypeCast
	case model.PoolPaysBenefitOtherUserTrTypeID.Int64():
		return int(model.PoolPaysBenefitOtherUserTrTypeID.Int64())
	default:
		return 0

	}
}

type IncomesHistory interface {
	GetNewIncomes(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error)
	GetIncomesIsViewed(ctx context.Context, id int, coinCode string) (bool, error)
}

type incomesHistory struct {
	pool          *pgxpool.Pool
	slaveDBPool   *pgxpool.Pool
	coinValidator coinValidatorRepo.CoinValidatorRepository
}

func NewIncomesHistory(
	pool *pgxpool.Pool,
	slaveDBPool *pgxpool.Pool,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
) IncomesHistory {
	return &incomesHistory{
		pool:          pool,
		slaveDBPool:   slaveDBPool,
		coinValidator: coinValidator,
	}
}

func (m *incomesHistory) GetNewIncomes(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error) {
	coinID, ok := m.coinValidator.GetIdByCode(filter.CoinCode)
	if !ok {
		return nil, fmt.Errorf("failed to find coinId by code: %s", filter.CoinCode)
	}

	var filterValues = []any{filter.UserID, coinID}

	values, query, countSumQuery, err := m.getNewIncomesFilterQuery(filterValues, filter)
	if err != nil {
		return nil, fmt.Errorf("m.getNewIncomesFilterQuery: %w", err)
	}

	rows, err := m.pool.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("selectIncomesNew Query: %w internal error 0002-01", err)
	}

	defer rows.Close()

	incomesRes := &model.HistoryOutput{
		TotalCount:    0,
		IncomesSum:    utils.DecimalToPtr(decimal.NewFromInt(0)),
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       make([]*model.Income, 0, filter.Limit),
		Payouts:       nil,
		Wallets:       nil,
	}

	incomeComplexities, err := m.getLastIncomeComplexity(ctx, int32(filter.UserID), coinID, filter.CoinCode)
	if err != nil {
		return nil, fmt.Errorf("getLastIncomeComplexity: %w internal error 0002-01", err)
	}

	for rows.Next() {
		var date time.Time
		income := new(model.Income)

		if err = rows.Scan(
			&date,
			&income.Code,
			&income.Time,
			&income.Income,
			&income.HashRate,
		); err != nil {
			return nil, fmt.Errorf("selectIncomesNew Scan: %w internal error 0002-02", err)
		}

		if income.Code == int(model.MainCoinMiningPayoutTrTypeID) {
			if complexity, ok := incomeComplexities[income.Time]; ok {
				income.ChangePercent = complexity.ChangePercent
				income.Diff = int(complexity.Diff)
			}
		}

		income.Code = incomeTypeCast(int64(income.Code))

		incomesRes.Incomes = append(incomesRes.Incomes, income)
	}

	if len(values) > veryCoolNumber2 {
		values = values[:len(values)-veryCoolNumber2]
	}

	countRows, err := m.pool.Query(ctx, countSumQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("countSumQuery Query: %w", err)
	}

	defer countRows.Close()

	if countRows.Next() {
		if err = countRows.Scan(&incomesRes.TotalCount, &incomesRes.IncomesSum); err != nil {
			return nil, fmt.Errorf("countSumQuery Scan: %w", err)
		}
	}

	return incomesRes, nil
}

func (m *incomesHistory) GetIncomesIsViewed(ctx context.Context, id int, coinCode string) (bool, error) {
	if coinCode == "" {
		coinCode = coinBTC
	}
	coinId, ok := m.coinValidator.GetIdByCode(coinCode)
	if !ok {
		return false, fmt.Errorf("failed to find coinId by code: %s", coinCode)
	}
	var isViewed bool
	if err := m.pool.QueryRow(ctx, selectIncomesIsViewed, id, coinId).Scan(&isViewed); err != nil {
		return false, fmt.Errorf("QueryRow: %w", err)
	}
	return isViewed, nil
}

func (m *incomesHistory) getNewIncomesFilterQuery(values []any, filter *model.HistoryInput) ([]any, string, string, error) {
	baseQuery := selectIncomesNew

	var filterQuery, paginationQuery string
	if filter.From != "" {
		from, err := time.Parse(time.DateOnly, filter.From)
		if err != nil {
			return nil, "", "", fmt.Errorf("filter.From: time.Parse: %w", err)
		}
		year, month, day := from.Date()
		values = append(values, fmt.Sprintf("%d-%d-%d", year, month, day))
		filterQuery += fmt.Sprintf(` AND date(t.created_at) >= $%d `, len(values))
	}

	if filter.To != "" {
		to, err := time.Parse(time.DateOnly, filter.To)
		if err != nil {
			return nil, "", "", fmt.Errorf("filter.To: time.Parse: %w", err)
		}
		year, month, day := to.Date()
		values = append(values, fmt.Sprintf("%d-%d-%d", year, month, day))
		filterQuery += fmt.Sprintf(` AND date(t.created_at) <= $%d `, len(values))
	}

	values = append(values, pq.Array(filter.TransactionTypesIDs))

	filterQuery += fmt.Sprintf(` AND t.type = ANY ($%d)`, len(values))

	countSumQuery := "SELECT COUNT(date), COALESCE(SUM(amount), 0) FROM (" + fmt.Sprintf(baseQuery, filterQuery, paginationQuery) + ") count_sum_table"

	values = append(values, filter.Limit)
	paginationQuery += fmt.Sprintf(` LIMIT $%d`, len(values))

	values = append(values, filter.Offset)
	paginationQuery += fmt.Sprintf(` OFFSET $%d`, len(values))

	return values, fmt.Sprintf(baseQuery, filterQuery, paginationQuery), countSumQuery, nil
}

func (m *incomesHistory) getLastIncomeComplexity(ctx context.Context, userID, coinID int32, coin string) (map[float64]*controller.Complexity, error) {
	incomeComplexities := make(map[float64]*controller.Complexity)

	if coin == "btc" || coin == "ltc" {
		rows, err := m.slaveDBPool.Query(ctx, selectComplexityIncomes, userID, coinID)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var iTimeList []float64

		for rows.Next() {
			var date time.Time
			var iTime float64

			if err := rows.Scan(
				&date,
				&iTime,
			); err != nil {
				return nil, err
			}

			iTimeList = append(iTimeList, iTime)
		}

		incomeComplexities = controller.GetComplexityByTimeList(coin, iTimeList)
	}

	return incomeComplexities, nil
}
