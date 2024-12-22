package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type LimitPayouts interface {
	GetRedisBlockStatus(ctx context.Context, userID int64) (int, error)
	SetRedisBlockStatus(ctx context.Context, userID int64, status int, exp time.Duration) error
	GetMainUserId(ctx context.Context, userID int64) (int64, error)
	GetLimit(ctx context.Context, coinID int) (float64, error)
	GetUserPayoutsSum(ctx context.Context, userID int64, coinId int) (float64, error)
	SetUserNopay(ctx context.Context, userID int64) error
}

type limitPayouts struct {
	pool       *pgxpool.Pool
	rds        *redis.Client
	coinLimits map[string]float64
}

func NewLimitPayouts(pool *pgxpool.Pool, redisCli *redis.Client, limits map[string]float64) LimitPayouts {
	return &limitPayouts{
		pool:       pool,
		rds:        redisCli,
		coinLimits: limits,
	}
}

func (r *limitPayouts) GetRedisBlockStatus(ctx context.Context, userID int64) (int, error) {
	key := fmt.Sprintf("payBlockStatus_%d", userID)
	val, err := r.rds.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("redis: %w", err)
	} else {
		if s, err := strconv.Atoi(val); err == nil {
			return s, nil
		} else {
			return 0, fmt.Errorf("string convert: %w", err)
		}
	}
}

func (r *limitPayouts) SetRedisBlockStatus(ctx context.Context, userID int64, status int, exp time.Duration) error {
	key := fmt.Sprintf("payBlockStatus_%d", userID)
	err := r.rds.Set(ctx, key, status, exp).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	return nil
}

func (r *limitPayouts) GetMainUserId(ctx context.Context, userID int64) (int64, error) {

	query := `
      SELECT COALESCE(parent_id, 0) FROM emcd.users WHERE id = $1
`
	var parentId int64
	row := r.pool.QueryRow(ctx, query, userID)
	err := row.Scan(&parentId)
	if err != nil {
		return 0, fmt.Errorf("repo get user parent_id: %w", err)
	}

	return parentId, nil

}

func (r *limitPayouts) SetUserNopay(ctx context.Context, userID int64) error {
	query := `
		UPDATE emcd.users SET nopay=true WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("postgres set nopay: %w", err)
	}

	return nil
}

func (r *limitPayouts) GetLimit(ctx context.Context, coinID int) (float64, error) {

	code, err := r.getCoinCodeById(ctx, coinID)
	if err != nil {
		return 0, err
	}

	limit, ok := r.coinLimits[code]
	if !ok {
		return 0, nil
	}

	return limit, nil
}

func (r *limitPayouts) getCoinCodeById(ctx context.Context, coinID int) (string, error) {
	query := `SELECT code FROM emcd.coins WHERE id=$1`
	var code string
	row := r.pool.QueryRow(ctx, query, coinID)
	err := row.Scan(&code)
	if err != nil {
		return "", fmt.Errorf("get coin by id: %w", err)
	} else {
		return code, nil
	}
}

func (r *limitPayouts) GetUserPayoutsSum(ctx context.Context, userID int64, coinId int) (float64, error) {

	query := `
    SELECT ROUND(COALESCE(SUM(t.amount), 0), 8) FROM emcd.transactions t
INNER JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
WHERE ua.user_id = $1 
  and t.coin_id = $2 
  and t.type = 30 
  and ua.account_type_id = 1 
  --and t.hash IS NOT NULL 
  and (t.hash != 'err' OR t.hash IS NULL)
  and t.created_at >= NOW() - INTERVAL '1 DAY'
`
	var result float64
	row := r.pool.QueryRow(ctx, query, userID, coinId)
	err := row.Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("scan row: %w", err)
	}

	return result, nil
}
