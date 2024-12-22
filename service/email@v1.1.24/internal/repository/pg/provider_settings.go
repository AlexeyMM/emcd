package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

type ProvideSettingsStore struct {
	pool *pgxpool.Pool
}

func (s *ProvideSettingsStore) Create(ctx context.Context, setting model.Setting) error {
	const createProviderSettingSQL = `
INSERT INTO provider_settings (whitelabel_id, providers, created_at, updated_at)
 VALUES($1, $2, $3, $4)
`
	_, err := s.pool.Exec(ctx, createProviderSettingSQL,
		setting.WhiteLabelID,
		&setting.Providers,
		setting.CreatedAt,
		setting.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("exec createProviderSettingSQL: %w", err)
	}
	return nil
}

func (s *ProvideSettingsStore) Get(ctx context.Context, whitelabelID uuid.UUID) (model.Setting, error) {
	const selectProviderSettingSQL = `
SELECT whitelabel_id, providers, created_at, updated_at 
  FROM provider_settings
 WHERE whitelabel_id = $1
`
	var result model.Setting
	err := s.pool.
		QueryRow(ctx, selectProviderSettingSQL, whitelabelID).
		Scan(&result.WhiteLabelID, &result.Providers, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, repository.ErrNotFound
		}
		return result, fmt.Errorf("exec selectProviderSettingSQL: %w", err)
	}
	return result, nil
}

func (s *ProvideSettingsStore) Update(ctx context.Context, setting model.Setting) error {
	const updateProviderSettingSQL = `
UPDATE provider_settings 
   SET providers = $1, 
       created_at = $2, 
       updated_at = $3 
 WHERE whitelabel_id = $4
`
	_, err := s.pool.Exec(
		ctx,
		updateProviderSettingSQL,
		setting.Providers,
		setting.CreatedAt,
		setting.UpdatedAt,
		setting.WhiteLabelID,
	)
	if err != nil {
		return fmt.Errorf("exec updateProviderSettingSQL: %w", err)
	}
	return nil
}

func (s *ProvideSettingsStore) List(ctx context.Context, p repository.Pagination) ([]model.Setting, error) {
	const listSettingsSQL = `
SELECT whitelabel_id, providers, created_at, updated_at
  FROM provider_settings
 ORDER BY whitelabel_id
 LIMIT $1 OFFSET $2
`
	rows, err := s.pool.Query(ctx, listSettingsSQL, p.Size, p.Offset())
	if err != nil {
		return nil, fmt.Errorf("exect listSettingsSQL: %w", err)
	}
	defer rows.Close()
	result := make([]model.Setting, 0, p.Size)
	var setting model.Setting
	for rows.Next() {
		err = rows.Scan(&setting.WhiteLabelID, &setting.Providers, &setting.CreatedAt, &setting.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan listSettingsSQL: %w", err)
		}
		result = append(result, setting)
	}
	return result, nil
}

func NewProvideSettingsStore(pool *pgxpool.Pool) *ProvideSettingsStore {
	return &ProvideSettingsStore{
		pool: pool,
	}
}
