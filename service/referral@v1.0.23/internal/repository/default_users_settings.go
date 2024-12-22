package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

const (
	approximatelyReferralCount = 100
)

type DefaultUserSettings struct {
	db *pgxpool.Pool
}

func NewDefaultUserSettings(db *pgxpool.Pool) *DefaultUserSettings {
	return &DefaultUserSettings{
		db: db,
	}
}

func (r *DefaultUserSettings) CreateUsersSettings(ctx context.Context, settings []model.ReferralSettings) error {
	createQuery := `
INSERT INTO 
    referral.default_settings_referrals(referral_id, product, coin, fee, referral_fee)
VALUES ($1, $2, $3, $4, $5)`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction CreateUsersSettings: %w", err)
	}

	batch := &pgx.Batch{}
	for _, user := range settings {
		for _, preference := range user.Preferences {
			batch.Queue(createQuery,
				user.ReferralUUID,
				preference.Product,
				preference.Coin,
				preference.Fee,
				preference.ReferralFee,
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

func (r *DefaultUserSettings) UpdateUsersSettings(
	ctx context.Context,
	settings []model.ReferralSettings,
	updateMode model.UpdateMode,
	usersWithPromoCodes map[string]map[string][]string) error {

	updateQueryReferralSettings := `
UPDATE referral.default_settings_referrals
	SET fee          = COALESCE($1, fee),
		referral_fee = COALESCE($2, referral_fee),
		updated_at = NOW()
WHERE referral_id = $3 AND product = $4 AND coin = $5
`

	updateCommission := `
UPDATE referral.referrals 
	SET fee          = COALESCE($1, fee),
		referral_fee = COALESCE($2, referral_fee)
WHERE referral_id = $3 AND product = $4 AND coin = $5
`

	updateCommissionForUsers := `
UPDATE referral.referrals 
	SET fee          = COALESCE($1, fee),
		referral_fee = COALESCE($2, referral_fee)
WHERE referral_id = $3 AND product = $4 AND coin = $5 AND NOT (user_id = ANY($6))
`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("open tx in UpdateUsersSettings: %w", err)
	}

	batch := &pgx.Batch{}
	for _, user := range settings {
		for _, preference := range user.Preferences {
			var referralFee sql.NullFloat64
			if preference.ReferralFee > 0 {
				referralFee.Float64 = preference.ReferralFee
				referralFee.Valid = true
			}

			var fee sql.NullFloat64
			if preference.Fee > 0 {
				fee.Float64 = preference.Fee
				fee.Valid = true
			}

			batch.Queue(updateQueryReferralSettings,
				fee,
				referralFee,
				user.ReferralUUID,
				preference.Product,
				preference.Coin,
			)

			switch updateMode {
			case model.UpdateModeForceAll:
				batch.Queue(updateCommission,
					fee,
					referralFee,
					user.ReferralUUID,
					preference.Product,
					preference.Coin,
				)
			case model.UpdateModeAll:
				if preference.Product != "mining" {
					continue
				}

				// если нет реферала в списке то обновляем всем пользователям
				coinAndUser, ok := usersWithPromoCodes[user.ReferralUUID.String()]
				if !ok {
					batch.Queue(updateCommission,
						fee,
						referralFee,
						user.ReferralUUID,
						preference.Product,
						preference.Coin,
					)
					continue
				}

				// если нет монеты то обновляем тоже всем
				users, ok := coinAndUser[preference.Coin]
				if !ok {
					batch.Queue(updateCommission,
						fee,
						referralFee,
						user.ReferralUUID,
						preference.Product,
						preference.Coin,
					)
					continue
				}

				// обновляем всем только не тем у кого промокод на эту монету
				batch.Queue(updateCommissionForUsers,
					fee,
					referralFee,
					user.ReferralUUID,
					preference.Product,
					preference.Coin,
					users,
				)
			}

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

func (r *DefaultUserSettings) GetUsersSettings(
	ctx context.Context,
	users []uuid.UUID,
	products, coins []string) (map[uuid.UUID]model.ReferralSettings, error) {
	args := []any{users}
	queryGet := `
SELECT referral_id, product, coin, fee, referral_fee, created_at, updated_at 
FROM referral.default_settings_referrals
WHERE referral_id = ANY($1)`

	num := 2
	if len(products) > 0 {
		queryGet += fmt.Sprintf(" AND product = ANY($%d)", num)
		args = append(args, products)
		num++
	}

	if len(coins) > 0 {
		queryGet += fmt.Sprintf(" AND coin = ANY($%d)", num)
		args = append(args, coins)
	}

	queryGet += " ORDER BY created_at ASC"

	rows, err := r.db.Query(ctx, queryGet, args...)
	if err != nil {
		return nil, fmt.Errorf("queryGet: %w", err)
	}
	defer rows.Close()

	list := make(map[uuid.UUID]model.ReferralSettings, approximatelyReferralCount)
	for rows.Next() {
		var (
			uuidRaw     sql.NullString
			product     sql.NullString
			coin        sql.NullString
			fee         sql.NullFloat64
			referralFee sql.NullFloat64
			createdAt   sql.NullTime
			updatedAt   sql.NullTime
		)
		err = rows.Scan(
			&uuidRaw,
			&product,
			&coin,
			&fee,
			&referralFee,
			&createdAt,
			&updatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		userUUID, err := uuid.Parse(uuidRaw.String)
		if err != nil {
			return nil, fmt.Errorf("parse uuid: %w", err)
		}

		user, ok := list[userUUID]
		if !ok {
			list[userUUID] = model.ReferralSettings{
				ReferralUUID: userUUID,
				Preferences: []model.ReferralPreference{
					{
						Product:     product.String,
						Coin:        coin.String,
						Fee:         fee.Float64,
						ReferralFee: referralFee.Float64,
						CreatedAt:   createdAt.Time,
						UpdatedAt:   updatedAt.Time,
					},
				},
			}
			continue
		}
		user.Preferences = append(user.Preferences, model.ReferralPreference{
			Product:     product.String,
			Coin:        coin.String,
			Fee:         fee.Float64,
			ReferralFee: referralFee.Float64,
			CreatedAt:   createdAt.Time,
			UpdatedAt:   updatedAt.Time,
		})

		list[userUUID] = user
	}

	return list, nil
}

func (r *DefaultUserSettings) GetUserUUIDsByReferralsUUID(ctx context.Context, settings []model.ReferralSettings) (map[uuid.UUID][]string, error) {
	getUsersByRefID := `
				SELECT distinct (user_id), referral_id
				FROM referral.referrals
				WHERE referral_id = ANY($1)
				`

	referralsUUIDs := make([]string, 0, len(settings))
	for _, set := range settings {
		referralsUUIDs = append(referralsUUIDs, set.ReferralUUID.String())
	}

	rows, err := r.db.Query(ctx, getUsersByRefID, referralsUUIDs)
	if err != nil {
		return nil, fmt.Errorf("getUsersByRefID: %w", err)
	}
	defer rows.Close()

	list := make(map[uuid.UUID][]string, approximatelyReferralCount)
	for rows.Next() {
		var (
			rawUUID         sql.NullString
			rawReferralUUID sql.NullString
		)
		err = rows.Scan(&rawUUID, &rawReferralUUID)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		userUUID, err := uuid.Parse(rawUUID.String)
		if err != nil {
			return nil, fmt.Errorf("parse uuid: %w", err)
		}
		referralUUID, err := uuid.Parse(rawReferralUUID.String)
		if err != nil {
			return nil, fmt.Errorf("parse referral uuid: %w", err)
		}

		referralUsersList, ok := list[referralUUID]
		if !ok {
			list[referralUUID] = []string{userUUID.String()}
			continue
		}

		referralUsersList = append(referralUsersList, userUUID.String())
		list[referralUUID] = referralUsersList
	}
	return list, nil
}
