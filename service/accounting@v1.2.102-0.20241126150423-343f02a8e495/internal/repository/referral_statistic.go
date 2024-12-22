package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"code.emcdtech.com/emcd/service/accounting/model"
)

const (
	getReferralsIncome = `
SELECT coin_id, SUM(amount) FROM emcd.transactions WHERE receiver_account_id=ANY($1) and type=ANY($2) AND DATE(created_at) >= $3 AND DATE(created_at) < $4 GROUP BY coin_id
    `
)

type ReferralStatistic interface {
	GetReferralsStatistic(ctx context.Context, req *model.ReferralsStatisticInput) (*model.ReferralsStatisticOutput, error)
}

type referralStatistic struct {
	pool *pgxpool.Pool
}

func NewReferralStatistic(pool *pgxpool.Pool) ReferralStatistic {
	return &referralStatistic{
		pool: pool,
	}
}

func (r *referralStatistic) GetReferralsStatistic(ctx context.Context, req *model.ReferralsStatisticInput) (*model.ReferralsStatisticOutput, error) {

	if len(req.AccountIDs) == 0 {
		return nil, fmt.Errorf("getReferralsStatistic: empty or nil list of account ids")
	}

	todayStart := time.Now().UTC().Truncate(24 * time.Hour)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	thisMonthStart := todayStart.AddDate(0, 0, -todayStart.Day()+1)

	yesterday, err := r.getReferralsStatisticForInterval(ctx, req.AccountIDs, req.TransactionTypesIDs, yesterdayStart, todayStart)
	if err != nil {
		return nil, fmt.Errorf("getReferralsStatistic (yesterday): %w: %s", err, Internal)

	}
	thisMonth, err := r.getReferralsStatisticForInterval(ctx, req.AccountIDs, req.TransactionTypesIDs, thisMonthStart, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("getReferralsStatistic (for this month): %w: %s", err, Internal)

	}

	return &model.ReferralsStatisticOutput{
		Income: model.AggregatedReferralIncome{
			Yesterday: yesterday,
			ThisMonth: thisMonth,
		},
	}, nil
}

func (r *referralStatistic) getReferralsStatisticForInterval(ctx context.Context, accountIDs, transactionTypesIDs []int64, from, to time.Time) ([]*model.ReferralIncome, error) {
	var sqlParams []any
	var query string
	sqlParams = append(sqlParams, pq.Array(accountIDs), pq.Array(transactionTypesIDs), from, to)
	query = getReferralsIncome

	var res []*model.ReferralIncome

	var amount float64
	var coinID int64

	rows, err := r.pool.Query(ctx, query, sqlParams...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getReferralsStatistic query row: %w: %s", err, Internal)
	}
	defer rows.Close()
	for rows.Next() {

		if err := rows.Scan(&coinID, &amount); err != nil {
			return nil, fmt.Errorf("getReferralsStatistic scanning row: %w: %s", err, Internal)
		}
		res = append(res, &model.ReferralIncome{CoinID: coinID, Amount: amount})
	}
	return res, nil
}
