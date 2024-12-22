package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/endpoint/internal/encryptor"
	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var secretNotFoundError = errors.New("secret not found")

type Secret struct {
	pool      *pgxpool.Pool
	encryptor encryptor.Encryptor
}

func NewSecret(pool *pgxpool.Pool, encryptor encryptor.Encryptor) *Secret {
	return &Secret{
		pool:      pool,
		encryptor: encryptor,
	}
}

func (s *Secret) Add(ctx context.Context, secret *model.Secret) error {
	encrypted, err := s.encryptor.Encrypt(secret.ApiSecret.String())
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	sql, args, err := squirrel.
		Insert("endpoint.secrets").
		Columns(
			"api_key",
			"api_secret",
			"client_id",
			"is_active",
			"created_at",
			"last_used").
		Values(
			secret.ApiKey,
			encrypted,
			secret.ClientID,
			secret.IsActive,
			secret.CreatedAt.UTC(),
			secret.LastUsed.UTC()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (s *Secret) Find(ctx context.Context, filter *model.SecretFilter) ([]*model.Secret, error) {
	return s.find(ctx, filter)
}

func (s *Secret) FindOne(ctx context.Context, filter *model.SecretFilter) (*model.Secret, error) {
	secrets, err := s.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(secrets) == 1 {
		return secrets[0], nil
	} else if len(secrets) > 1 {
		return nil, fmt.Errorf("unexpected number of secrets: %d", len(secrets))
	}
	return nil, secretNotFoundError
}

func (s *Secret) find(ctx context.Context, filter *model.SecretFilter) ([]*model.Secret, error) {
	query := squirrel.Select(
		"api_key",
		"api_secret",
		"client_id",
		"is_active",
		"created_at",
		"last_used").
		From("endpoint.secrets").
		PlaceholderFormat(squirrel.Dollar)

	query = newSecretFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var secrets []*model.Secret
	for rows.Next() {
		var (
			secret                     model.Secret
			encryptedSecret, apiSecret string
			apiSecretUid               uuid.UUID
		)
		err = rows.Scan(&secret.ApiKey, &encryptedSecret, &secret.ClientID, &secret.IsActive, &secret.CreatedAt, &secret.LastUsed)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		apiSecret, err = s.encryptor.Decrypt(encryptedSecret)
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}

		apiSecretUid, err = uuid.Parse(apiSecret)
		if err != nil {
			return nil, fmt.Errorf("parse api secret: %w", err)
		}

		secret.ApiSecret = apiSecretUid

		secrets = append(secrets, &secret)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return secrets, nil
}

func (s *Secret) Update(ctx context.Context, secret *model.Secret, filter *model.SecretFilter, partial *model.SecretPartial) error {
	query := squirrel.
		Update("endpoint.secrets").
		PlaceholderFormat(squirrel.Dollar)

	query = newSecretFilterSql(filter).applyToUpdateQuery(query)

	query = newSecretPartialSql(partial).applyToQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	secret.Update(partial)

	return nil
}
