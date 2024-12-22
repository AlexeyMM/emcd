// Package repository provides functions for managing data storage and retrieval.
package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

// DefaultWhitelabelSettings is an interface that defines the contract for working with default whitelabel settings.
// It provides methods for creating, updating, deleting, retrieving, and retrieving all default whitelabel settings.
// The methods take a context, various input parameters, and return the requested data or an common if any occurred.
//
//go:generate mockery --name=DefaultWhitelabelSettings --structname=MockDefaultWhitelabelSettings --outpkg=repository --output ./ --filename default_whitelabel_settings_mock.go
type DefaultWhitelabelSettings interface {
	Create(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error
	Update(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error
	Delete(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) error
	GetAllWithFilters(
		ctx context.Context,
		skip int32,
		take int32,
		filters map[string]string,
	) ([]*model.DefaultWhitelabelSettingsV2, int, error)
	GetAllWithoutPaginationWithFilters(
		ctx context.Context,
		filters map[string]string,
	) ([]*model.DefaultWhitelabelSettingsV2, error)
	GetV2(ctx context.Context, wlID uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error)
	GetV2ByCoin(ctx context.Context, product, coin string, wlID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)
	Get(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)
}

// defaultWhitelabelSettings is a type that represents the default whitelabel settings.
type defaultWhitelabelSettings struct {
	db *pgxpool.Pool
}

// NewDefaultWhitelabelSettings creates a new instance of DefaultWhitelabelSettings interface.
// It takes a *pgxpool.Pool as a parameter and returns a DefaultWhitelabelSettings object.
// The function constructs a defaultWhitelabelSettings object with the provided *pgx.Conn as the db field.
// It then returns the created object as a DefaultWhitelabelSettings interface.
func NewDefaultWhitelabelSettings(db *pgxpool.Pool) DefaultWhitelabelSettings {
	return &defaultWhitelabelSettings{db: db}
}

// Create inserts a new entry into the default_whitelabel_settings table.
// It takes a context and a DefaultWhitelabelSettings object as parameters.
// The function constructs an INSERT query using the provided object and executes it.
// If the execution encounters an common, it returns an common with the specific details.
// Otherwise, it returns nil.
func (s *defaultWhitelabelSettings) Create(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	query := `
				INSERT INTO referral.default_whitelabel_settings(product, whitelabel_id, coin, fee, referral_fee,whitelabel_fee)
				VALUES ($1, $2, $3, $4, $5, $6)
				`

	_, err := s.db.Exec(ctx, query,
		in.Product,
		in.WhitelabelID,
		in.Coin,
		in.Fee,
		in.ReferralFee,
		in.WhiteLabelFee)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Update updates the fee and referral fee of a specific row in the default_whitelabel_settings table.
// It takes a context and a DefaultWhitelabelSettings object as parameters.
// The function constructs an UPDATE query using the provided object and executes it.
// If the execution encounters an common, it returns an common with the specific details.
// Otherwise, it returns nil.
func (s *defaultWhitelabelSettings) Update(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	query := `
				UPDATE referral.default_whitelabel_settings
				SET fee          = $4,
					referral_fee = $5,
					whitelabel_fee = $6
				WHERE product = $1
				  AND whitelabel_id = $2
				  AND coin = $3
				`

	_, err := s.db.Exec(ctx, query,
		in.Product,
		in.WhitelabelID,
		in.Coin,
		in.Fee,
		in.ReferralFee,
		in.WhiteLabelFee)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// Delete removes an entry from the default_whitelabel_settings table based on the given parameters.
// It takes a context, product, coin, and whitelabelID as parameters.
// The function constructs a DELETE query using the provided parameters and executes it.
// If the execution encounters an common, it returns an common with the specific details.
// Otherwise, it returns nil.
func (s *defaultWhitelabelSettings) Delete(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) error {
	query := `
				DELETE
				FROM referral.default_whitelabel_settings
				WHERE product = $1
				  AND whitelabel_id = $2
				  AND coin = $3
				`

	_, err := s.db.Exec(ctx, query,
		product,
		whitelabelID,
		coin)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// GetAll retrieves a list of DefaultWhitelabelSettings objects from the default_whitelabel_settings table.
// It takes a context, a skip parameter indicating the number of items to skip, and a take parameter indicating the maximum number of items to retrieve.
// The function constructs a SELECT query using the provided parameters and executes it.
// If the execution encounters an common, it returns nil for the list and 0 for the total count with an common message.
// Otherwise, it retrieves the rows returned by the query and scans them into DefaultWhitelabelSettings objects.
// The function then appends each object to the list.
// After retrieving all the rows, it constructs a separate query to get the total count of rows in the default_whitelabel_settings table.
// Finally, it scans the total count and returns the list of objects, the total count, and nil for the common if successful.
// If there is an common while executing the count query, it returns nil for the list, 0 for the total count, and an common message.
func (s *defaultWhitelabelSettings) GetAllWithFilters(
	ctx context.Context,
	skip int32,
	take int32,
	filters map[string]string,
) ([]*model.DefaultWhitelabelSettingsV2, int, error) {
	filterQuery, filterArgs := s.getFilterQuery(filters)

	batch := pgx.Batch{}

	query := `
				SELECT product, whitelabel_id, coin, fee, referral_fee, created_at, whitelabel_fee
				FROM referral.default_whitelabel_settings %s ORDER BY created_at
				LIMIT $%d OFFSET $%d
				`
	query = fmt.Sprintf(query, filterQuery, len(filterArgs)+1, len(filterArgs)+2)

	log.Info(ctx, query)
	batch.Queue(query, append(filterArgs, take, skip)...)

	query = `
				SELECT COUNT(*)
				FROM referral.default_whitelabel_settings %s
				`
	query = fmt.Sprintf(query, filterQuery)
	batch.Queue(query, filterArgs...)

	res := s.db.SendBatch(ctx, &batch)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(ctx, "defaultWhitelabelSettings.GetAll: close batch: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("query %w", err)
	}

	defer rows.Close()

	list := make([]*model.DefaultWhitelabelSettingsV2, 0)

	for rows.Next() {
		o := new(model.DefaultWhitelabelSettingsV2)

		err = rows.Scan(
			&o.Product,
			&o.WhitelabelID,
			&o.Coin,
			&o.Fee,
			&o.ReferralFee,
			&o.CreatedAt,
			&o.WhiteLabelFee)

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

func (s *defaultWhitelabelSettings) GetAllWithoutPaginationWithFilters(
	ctx context.Context,
	filters map[string]string,
) ([]*model.DefaultWhitelabelSettingsV2, error) {
	filterQuery, filterArgs := s.getFilterQuery(filters)

	query := `SELECT product, whitelabel_id, coin, fee, referral_fee, created_at, whitelabel_fee
				FROM referral.default_whitelabel_settings %s ORDER BY created_at`
	query = fmt.Sprintf(query, filterQuery)

	rows, err := s.db.Query(ctx, query, filterArgs...)
	if err != nil {
		return nil, fmt.Errorf("queery: %w", err)
	}
	defer rows.Close()
	res := make([]*model.DefaultWhitelabelSettingsV2, 0)
	for rows.Next() {
		var ds model.DefaultWhitelabelSettingsV2
		err = rows.Scan(&ds.Product, &ds.WhitelabelID, &ds.Coin, &ds.Fee, &ds.ReferralFee, &ds.CreatedAt, &ds.WhiteLabelFee)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, &ds)
	}
	return res, nil
}

func (s *defaultWhitelabelSettings) GetV2(ctx context.Context, wlID uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error) {
	query := `SELECT product,whitelabel_id,coin,fee,referral_fee, created_at, whitelabel_fee FROM referral.default_whitelabel_settings WHERE whitelabel_id=$1`
	rows, err := s.db.Query(ctx, query, wlID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	res := make([]*model.DefaultWhitelabelSettingsV2, 0)
	for rows.Next() {
		var d model.DefaultWhitelabelSettingsV2
		err = rows.Scan(&d.Product, &d.WhitelabelID, &d.Coin, &d.Fee, &d.ReferralFee, &d.CreatedAt, &d.WhiteLabelFee)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, &d)
	}
	return res, nil
}

func (s *defaultWhitelabelSettings) GetV2ByCoin(
	ctx context.Context,
	product, coin string,
	wlID uuid.UUID,
) (*model.DefaultWhitelabelSettingsV2, error) {
	query := `SELECT product,whitelabel_id,coin,fee,referral_fee, created_at, whitelabel_fee FROM referral.default_whitelabel_settings WHERE whitelabel_id=$1 AND product=$2 AND coin=$3`
	var d model.DefaultWhitelabelSettingsV2
	err := s.db.QueryRow(ctx, query, wlID, product, coin).
		Scan(&d.Product, &d.WhitelabelID, &d.Coin, &d.Fee, &d.ReferralFee, &d.CreatedAt, &d.WhiteLabelFee)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	return &d, nil
}

// Get retrieves a default whitelabel settings entry from the default_whitelabel_settings table.
// It takes a context and three strings (product, coin, and whitelabelID) as parameters.
// The function constructs a SELECT query using the provided parameters and executes it.
// If a matching entry is found, the values are scanned into a DefaultWhitelabelSettings object.
// If the query encounters an common or no matching entry is found, the function returns an common with specific details.
// Otherwise, it returns the scanned DefaultWhitelabelSettings object and nil.
func (s *defaultWhitelabelSettings) Get(
	ctx context.Context,
	product string,
	coin string,
	whitelabelID uuid.UUID,
) (*model.DefaultWhitelabelSettingsV2, error) {
	query := `
				SELECT product, whitelabel_id, coin, fee, referral_fee, created_at, whitelabel_fee
				FROM referral.default_whitelabel_settings
				WHERE product = $1
				  AND whitelabel_id = $2
				  AND coin = $3
				`

	out := new(model.DefaultWhitelabelSettingsV2)

	err := s.db.QueryRow(ctx, query,
		product,
		whitelabelID,
		coin).Scan(
		&out.Product,
		&out.WhitelabelID,
		&out.Coin,
		&out.Fee,
		&out.ReferralFee,
		&out.CreatedAt,
		&out.WhiteLabelFee)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSettingsNotFound
		}

		return nil, fmt.Errorf("queryRow %w", err)
	}

	return out, nil
}

func (s *defaultWhitelabelSettings) getFilterQuery(filters map[string]string) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}
	res := "WHERE "
	filterVals := make([]interface{}, 0)
	i := 1
	for name, val := range filters {
		res += fmt.Sprintf("%s=$%d AND ", name, i)
		filterVals = append(filterVals, val)
		i++
	}
	res = res[:len(res)-4]
	return res, filterVals
}
