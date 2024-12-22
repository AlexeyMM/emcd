// Package repository provides functions for managing data storage and retrieval.
package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

// Referral interface represents a referral service.
//
//go:generate mockery --name=Referral --structname=MockReferral --outpkg=repository --output ./ --filename referral_mock.go
type Referral interface {
	Create(ctx context.Context, in *model.Referral) error
	Update(ctx context.Context, in *model.Referral) error
	Delete(ctx context.Context, userID uuid.UUID, product string, coin string) error
	Get(ctx context.Context, userID uuid.UUID, product string, coin string) (*model.Referral, error)
	List(ctx context.Context, userID uuid.UUID, skip int32, take int32) ([]*model.Referral, int, error)
	History(ctx context.Context, userID uuid.UUID, product string, coin string) ([]*model.Referral, error)
	CreateMultiple(ctx context.Context, rs []*model.Referral) error
	UpdateWithMultiplier(ctx context.Context, userID uuid.UUID, product string, coins []string, multiplier decimal.Decimal) error
	GetUserReferrals(ctx context.Context, userID uuid.UUID, skip, limit int32) ([]*model.UserReferral, int64, error)
	UpdateFee(ctx context.Context, userID uuid.UUID, product string, fees map[string]decimal.Decimal) error
	UpdateWithPromoCodeByCoin(ctx context.Context, cm *model.CoinMultiplier) error

	UpdateFeeByCoinAndProduct(ctx context.Context, userUUID []string, fees []model.SettingForCoinAndProduct) error
	UpdateReferralUUIDByUserUUID(ctx context.Context, userUUIDs []string, referralUUID string) error
}

// referral represents a referral service.
type referral struct {
	trx pgTx.PgxTransactor
	db  *pgxpool.Pool
}

// NewReferral creates a new instance of the Referral interface with the provided db connection.
// It returns a reference to the created Referral instance.
func NewReferral(trx pgTx.PgxTransactor, db *pgxpool.Pool) Referral {
	return &referral{
		trx: trx,
		db:  db,
	}
}

// Create inserts a new referral record into the database.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes a pointer to a model.Referral object, which contains the referral data to be inserted.
// It returns an common if any common occurs during the execution of the database query.
func (r *referral) Create(ctx context.Context, in *model.Referral) error {
	query := `
				INSERT INTO referral.referrals(user_id, product, coin, whitelabel_id, fee, whitelabel_fee, referral_fee, referral_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				`
	var err error

	_, err = r.trx.Runner(ctx).Exec(ctx, query,
		in.UserID,
		in.Product,
		in.Coin,
		in.WhitelabelID,
		in.Fee,
		in.WhitelabelFee,
		in.ReferralFee,
		in.ReferralID)

	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Update updates an existing referral record in the database.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes a pointer to a model.Referral object, which contains the updated referral data.
// It returns an common if any common occurs during the execution of the database query.
func (r *referral) Update(ctx context.Context, in *model.Referral) error {
	query := `
				UPDATE referral.referrals
				SET fee            = $5,
				    whitelabel_fee = $6,
					referral_fee   = $7,
					referral_id    = $8
				WHERE user_id = $1
				  AND product = $2
				  AND coin = $3
				  AND whitelabel_id = $4
				`

	var err error

	_, err = r.trx.Runner(ctx).Exec(ctx, query,
		in.UserID,
		in.Product,
		in.Coin,
		in.WhitelabelID,
		in.Fee,
		in.WhitelabelFee,
		in.ReferralFee,
		in.ReferralID)

	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Delete removes a referral record from the database based on the specified user ID, product, and coin.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes the user ID, product, and coin of the referral record to be deleted.
// It returns an common if any common occurs during the execution of the database query.
func (r *referral) Delete(ctx context.Context, userID uuid.UUID, product string, coin string) error {
	query := `
				DELETE
				FROM referral.referrals
				WHERE user_id = $1
				  AND product = $2
				  AND coin = $3
				`

	_, err := r.trx.Runner(ctx).Exec(ctx, query,
		userID,
		product,
		coin)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Get retrieves a referral record from the database based on the provided userID, product, and coin.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes three strings: userID, product, and coin, which are used as the parameters for the database query.
// It returns a pointer to a model.Referral object and an common.
// The model.Referral object contains the retrieved referral data.
// The common value indicates whether any common occurred during the execution of the database query.
func (r *referral) Get(ctx context.Context, userID uuid.UUID, product string, coin string) (*model.Referral, error) {
	query := `
				SELECT user_id,
					   product,
					   coin,
					   whitelabel_id,
					   fee,
					   whitelabel_fee,
					   referral_fee,
					   referral_id,
					   created_at
				FROM referral.referrals
				WHERE user_id = $1
				  AND product = $2
				  AND coin = $3
				`
	var err error

	out := new(model.Referral)

	err = r.trx.Runner(ctx).QueryRow(ctx, query,
		userID,
		product,
		coin).Scan(
		&out.UserID,
		&out.Product,
		&out.Coin,
		&out.WhitelabelID,
		&out.Fee,
		&out.WhitelabelFee,
		&out.ReferralFee,
		&out.ReferralID,
		&out.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSettingsNotFound
		}

		return nil, err
	}

	return out, nil
}

func (r *referral) GetUserReferrals(ctx context.Context, userID uuid.UUID, skip, limit int32) ([]*model.UserReferral, int64, error) {
	b := pgx.Batch{}
	getUsersByRefID := `
				SELECT distinct (user_id)
				FROM referral.referrals
				WHERE referral_id = $1
				LIMIT $2 OFFSET $3
				`
	b.Queue(getUsersByRefID, userID, limit, skip)

	const countQuery = `
SELECT COUNT(DISTINCT user_id)
	FROM referral.referrals
	WHERE referral_id = $1`
	b.Queue(countQuery, userID)

	res := r.trx.Runner(ctx).SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error(ctx, "GetUserReferrals: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("GetUserReferrals: query getUsersByRefID: %w", err)
	}
	defer rows.Close()

	users := make([]*model.UserReferral, 0)
	for rows.Next() {
		var s model.UserReferral
		err = rows.Scan(&s.UserID)
		if err != nil {
			return nil, 0, errors.Wrap(err, "GetUserReferrals db scan")
		}

		users = append(users, &s)
	}

	var count int64
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("GetUserReferrals: query countQuery: %w", err)
	}

	return users, count, nil
}

// List retrieves a list of referrals for a specific user from the database.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes the userID of the user for whom to retrieve the referrals.
// The skip parameter specifies the number of records to skip before retrieving the results.
// The take parameter specifies the maximum number of records to retrieve.
// It returns a slice of model.Referral objects representing the retrieved referrals.
// The second return value is the total count of referrals for the user.
// It also returns an common if any common occurs during the execution of the database query.
func (r *referral) List(ctx context.Context, userID uuid.UUID, skip int32, take int32) ([]*model.Referral, int, error) {
	batch := pgx.Batch{}

	query := `
				SELECT user_id,
					   product,
					   coin,
					   whitelabel_id,
					   fee,
					   whitelabel_fee,
					   referral_fee,
					   referral_id,
					   created_at
				FROM referral.referrals
				WHERE user_id = $1
				LIMIT $2 OFFSET $3
				`

	batch.Queue(query, userID, take, skip)

	query = `
				SELECT COUNT(*)
				FROM referral.referrals
				WHERE user_id = $1
				`

	batch.Queue(query, userID)

	res := r.trx.Runner(ctx).SendBatch(ctx, &batch)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(ctx, "referral.List: close batch: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("query %w", err)
	}

	defer rows.Close()

	list := make([]*model.Referral, 0)

	for rows.Next() {
		o := new(model.Referral)

		err = rows.Scan(
			&o.UserID,
			&o.Product,
			&o.Coin,
			&o.WhitelabelID,
			&o.Fee,
			&o.WhitelabelFee,
			&o.ReferralFee,
			&o.ReferralID,
			&o.CreatedAt)

		if err != nil {
			return nil, 0, fmt.Errorf("scan %w", err)
		}

		list = append(list, o)
	}

	var totalCount int

	err = res.QueryRow().Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("queryRow %w", err)
	}

	return list, totalCount, nil
}

// History retrieves the history of referrals for a specific user, product, and coin from the referral.referrals_logs table in the database.
// It takes a context, which carries deadlines, cancellation signals, and other values across API boundaries.
// It also takes the userID, product, and coin as parameters to filter the results.
// It returns a slice of model.Referral objects representing the referral history and an common if any common occurs during the execution of the database query.
func (r *referral) History(ctx context.Context, userID uuid.UUID, product string, coin string) ([]*model.Referral, error) {
	query := `
				SELECT data ->> 'user_id'                   AS user_id,
				       data ->> 'product'                   AS product,
					   data ->> 'coin'                      AS coin,
					   data ->> 'whitelabel_id'             AS whitelabel_id,
					   data ->> 'referral_id'               AS referral_id,
					   (data ->> 'fee')::decimal            AS fee,
					   (data ->> 'whitelabel_fee')::decimal AS whitelabel_fee,
					   (data ->> 'referral_fee')::decimal   AS referral_fee,
					   (data ->> 'created_at')::timestamp   AS created_at
				FROM referral.referrals_logs,
					 jsonb_array_elements(history) AS data
				WHERE user_id = $1
				  AND product = $2
				  AND coin = $3
				`

	var history []*model.Referral

	rows, err := r.trx.Runner(ctx).Query(ctx, query, userID, product, coin)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		o := new(model.Referral)

		err = rows.Scan(
			&o.UserID,
			&o.Product,
			&o.Coin,
			&o.WhitelabelID,
			&o.ReferralID,
			&o.Fee,
			&o.WhitelabelFee,
			&o.ReferralFee,
			&o.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan %w", err)
		}

		history = append(history, o)
	}

	return history, nil
}

const numberOfParamsInReferral = 9

func (r *referral) CreateMultiple(ctx context.Context, rs []*model.Referral) error {
	query := `INSERT INTO referral.referrals (user_id, product, coin, whitelabel_id, fee, whitelabel_fee, referral_fee, referral_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction CreateUsersSettings: %w", err)
	}
	batch := &pgx.Batch{}
	for _, rT := range rs {
		batch.Queue(query,
			rT.UserID,
			rT.Product,
			rT.Coin,
			rT.WhitelabelID,
			rT.Fee,
			rT.WhitelabelFee,
			rT.ReferralFee,
			rT.ReferralID,
			rT.CreatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("exec queryInsert: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *referral) UpdateWithMultiplier(
	ctx context.Context,
	userID uuid.UUID,
	product string,
	coins []string,
	multiplier decimal.Decimal,
) error {
	query := `UPDATE referral.referrals SET fee=ROUND($1*fee,4) WHERE user_id=$2 AND product=$3 AND coin=ANY($4)`
	_, err := r.trx.Runner(ctx).Exec(ctx, query, multiplier, userID, product, coins)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *referral) UpdateFee(ctx context.Context, userID uuid.UUID, product string, fees map[string]decimal.Decimal) error {
	tx, err := r.trx.Runner(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `UPDATE referral.referrals SET fee=$1 WHERE user_id=$2 AND product=$3 AND coin=$4`
	for coin, fee := range fees {
		_, err = tx.Exec(ctx, query, fee, userID, product, coin)
		if err != nil {
			return fmt.Errorf("exec: fee: %s. coin: %s. %w", fee.String(), coin, err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (r *referral) UpdateWithPromoCodeByCoin(ctx context.Context, cm *model.CoinMultiplier) error {
	err := r.trx.WithinTransaction(ctx, func(ctx context.Context) error {
		query := `UPDATE referral.referrals SET fee=fee*$1, referral_fee=referral_fee*$2 WHERE user_id=$3 AND product=$4 AND coin=$5`
		_, inTxErr := r.trx.Runner(ctx).Exec(ctx, query, cm.FeeMultiplier, cm.RefFeeMultiplier, cm.UserID, cm.Product,
			cm.Coin)
		if inTxErr != nil {
			return fmt.Errorf("update referrals: exec: %w", inTxErr)
		}
		query = `INSERT INTO referral.promocodes_updates (user_id,action_id,coin,created_at) VALUES ($1,$2,$3,$4)`
		_, inTxErr = r.trx.Runner(ctx).Exec(ctx, query, cm.UserID, cm.ActionID, cm.Coin, cm.CreatedAt)
		if inTxErr != nil {
			return fmt.Errorf("insert promocodes_updates: exec: %w", inTxErr)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("trx.WithinTransaction: %w", err)
	}
	return nil
}

func (r *referral) UpdateFeeByCoinAndProduct(ctx context.Context, usersUUIDs []string, fees []model.SettingForCoinAndProduct) error {
	updateQueryReferral := `
UPDATE referral.referrals 
SET fee=$1
WHERE user_id = $2 AND product = $3 AND coin = $4
`
	tx, err := r.trx.Runner(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("open tx in UpdateFeeByCoinAndProduct: %w", err)
	}

	batch := &pgx.Batch{}
	for _, userUUID := range usersUUIDs {
		for _, fee := range fees {
			batch.Queue(updateQueryReferral,
				fee.Fee,
				userUUID,
				fee.Product,
				fee.Coin,
			)
		}
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("exec queryInsert: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *referral) UpdateReferralUUIDByUserUUID(ctx context.Context, userUUIDs []string, referralUUID string) error {
	updateQueryReferralUUID := `UPDATE referral.referrals SET referral_id=$1 WHERE user_id=$2`

	tx, err := r.trx.Runner(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("open tx in UpdateReferralUUIDByUserUUID: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	batch := &pgx.Batch{}
	for _, userUUID := range userUUIDs {
		batch.Queue(updateQueryReferralUUID,
			referralUUID,
			userUUID,
		)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("exec queryInsert: %w", err)
	}

	err = br.Close()
	if err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
