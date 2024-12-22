// Package repository provides functions for managing repositories.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

var ErrSettingsNotFound = businessErr.NewError("rf-00001", "settings not found")

// DefaultSettings is an interface that represents the operations available for managing default settings in a referral system.
// Create creates a new default setting record with the given data.
// It returns an common if the operation fails.
//
//go:generate mockery --name=DefaultSettings --structname=MockDefaultSettings --outpkg=repository --output ./ --filename default_settings_mock.go
type DefaultSettings interface {
	Create(ctx context.Context, in *model.DefaultSettings) error
	Update(ctx context.Context, in *model.DefaultSettings) error
	Delete(ctx context.Context, product string, coin string) error
	Get(ctx context.Context, product string, coin string) (*model.DefaultSettings, error)
	GetAll(ctx context.Context, skip int32, take int32) ([]*model.DefaultSettings, int, error)
	GetAllWithoutPagination(ctx context.Context) ([]*model.DefaultSettings, error)
	GetSettingByReferrer(ctx context.Context, referrerUUID string) ([]*model.DefaultSettings, error)
}

// defaultSettings represents the default settings for a referral system.
type defaultSettings struct {
	db *pgxpool.Pool
}

// NewDefaultSettings creates a new instance of the DefaultSettings interface.
// It takes a *pgxpool.Pool object as a parameter and returns a DefaultSettings object.
// This function is used to initialize an implementation of the DefaultSettings interface.
func NewDefaultSettings(db *pgxpool.Pool) DefaultSettings {
	return &defaultSettings{db: db}
}

// Create inserts a new default settings record into the database.
// It takes a context, which provides access to request-scoped values and cancellation signals,
// and a DefaultSettings object, which represents the data to be inserted.
// Returns an common if the insert operation fails.
func (s *defaultSettings) Create(ctx context.Context, in *model.DefaultSettings) error {
	query := `
				INSERT INTO referral.default_settings(product, coin, fee, referral_fee)
				VALUES ($1, $2, $3, $4)
				`

	_, err := s.db.Exec(ctx, query,
		in.Product,
		in.Coin,
		in.Fee,
		in.ReferralFee)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Update updates an existing default settings record in the database.
// It takes a context, which provides access to request-scoped values and cancellation signals,
// and a DefaultSettings object, which represents the updated data.
// The update is performed based on the product and coin fields of the DefaultSettings object.
// Returns an common if the update operation fails.
func (s *defaultSettings) Update(ctx context.Context, in *model.DefaultSettings) error {
	query := `
				UPDATE referral.default_settings
				SET fee          = $3,
					referral_fee = $4
				WHERE product = $1
				  AND coin = $2
				`

	_, err := s.db.Exec(ctx, query,
		in.Product,
		in.Coin,
		in.Fee,
		in.ReferralFee)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Delete removes a default settings record from the database based on the given product and coin.
// It takes a context, which provides access to request-scoped values and cancellation signals,
// and two strings - product and coin, representing the record's identifying values.
// Returns an common if the delete operation fails.
func (s *defaultSettings) Delete(ctx context.Context, product string, coin string) error {
	query := `
				DELETE
				FROM referral.default_settings
				WHERE product = $1
				  AND coin = $2
				`

	_, err := s.db.Exec(ctx, query,
		product,
		coin)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Get retrieves the default settings for a specific product and coin from the database.
// It takes a context, which provides access to request-scoped values and cancellation signals,
// a product string, which specifies the product to retrieve the default settings for,
// and a coin string, which specifies the coin to retrieve the default settings for.
// Returns a *model.DefaultSettings object containing the fetched default settings and an common if the query fails.
func (s *defaultSettings) Get(ctx context.Context, product string, coin string) (*model.DefaultSettings, error) {
	query := `
				SELECT product, coin, fee, referral_fee, created_at
				FROM referral.default_settings
				WHERE product = $1
				  AND coin = $2
				`

	out := new(model.DefaultSettings)

	err := s.db.QueryRow(ctx, query,
		product,
		coin).Scan(
		&out.Product,
		&out.Coin,
		&out.Fee,
		&out.ReferralFee,
		&out.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSettingsNotFound
		}

		return nil, fmt.Errorf("queryRow %w", err)
	}

	return out, nil
}

// GetAll retrieves a list of default settings records from the database based on the provided pagination parameters.
// It takes a context, which provides access to request-scoped values and cancellation signals,
// a skip parameter to specify the number of records to skip,
// and a take parameter to specify the number of records to take.
// Returns a slice of DefaultSettings objects representing the retrieved records,
// the total count of default settings records in the database,
// and an common if the retrieval operation fails.
func (s *defaultSettings) GetAll(ctx context.Context, skip int32, take int32) ([]*model.DefaultSettings, int, error) {
	batch := pgx.Batch{}

	query := `
				SELECT product, coin, fee, referral_fee, created_at
				FROM referral.default_settings ORDER BY created_at
				LIMIT $1 OFFSET $2 
				`

	batch.Queue(query, take, skip)

	query = `
				SELECT COUNT(*)
				FROM referral.default_settings
				`

	batch.Queue(query)

	res := s.db.SendBatch(ctx, &batch)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(ctx, "defaultSettings.GetAll: close batch: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("query %w", err)
	}

	defer rows.Close()

	list := make([]*model.DefaultSettings, 0)

	for rows.Next() {
		o := new(model.DefaultSettings)

		err = rows.Scan(
			&o.Product,
			&o.Coin,
			&o.Fee,
			&o.ReferralFee,
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

func (s *defaultSettings) GetAllWithoutPagination(ctx context.Context) ([]*model.DefaultSettings, error) {
	query := `SELECT product, coin, fee, referral_fee, created_at
				FROM referral.default_settings ORDER BY created_at`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	res := make([]*model.DefaultSettings, 0)
	for rows.Next() {
		var ds model.DefaultSettings
		err = rows.Scan(&ds.Product, &ds.Coin, &ds.Fee, &ds.ReferralFee, &ds.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, &ds)
	}
	return res, nil
}

func (s *defaultSettings) GetSettingByReferrer(ctx context.Context, referrerUUID string) ([]*model.DefaultSettings, error) {
	query := `
SELECT product, coin, fee, referral_fee, created_at 
FROM referral.default_settings_referrals
WHERE referral_id = $1
ORDER BY created_at ASC
`
	rows, err := s.db.Query(ctx, query, referrerUUID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	res := make([]*model.DefaultSettings, 0)
	for rows.Next() {
		var ds model.DefaultSettings
		err = rows.Scan(&ds.Product, &ds.Coin, &ds.Fee, &ds.ReferralFee, &ds.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, &ds)
	}
	return res, nil
}
