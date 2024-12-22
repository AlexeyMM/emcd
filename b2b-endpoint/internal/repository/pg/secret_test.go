package pg

import (
	"context"
	"testing"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSecret_Add(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	secretRepo := NewSecret(db, encrypt)
	clientRepo := NewClient(db)

	cli := &model.Client{
		ID:   uuid.New(),
		Name: "Test Client",
	}
	err := clientRepo.Add(ctx, cli)
	require.NoError(t, err)

	// Пример данных для секрета
	secret := &model.Secret{
		ApiKey:    uuid.New(),
		ApiSecret: uuid.New(),
		ClientID:  cli.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Выполняем тестируемую функцию Add
	err = secretRepo.Add(ctx, secret)
	require.NoError(t, err)

	foundSecret, err := secretRepo.FindOne(ctx, &model.SecretFilter{
		ApiKey: &secret.ApiKey,
	})
	require.NoError(t, err)
	require.Equal(t, secret.ApiKey, foundSecret.ApiKey)
	require.Equal(t, secret.ApiSecret, foundSecret.ApiSecret)
	require.Equal(t, secret.ClientID, foundSecret.ClientID)
	require.Equal(t, secret.IsActive, foundSecret.IsActive)
	require.WithinDuration(t, secret.CreatedAt, foundSecret.CreatedAt, time.Second)
	require.WithinDuration(t, secret.LastUsed, foundSecret.LastUsed, time.Second)
}

func TestSecret_Find(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	secretRepo := NewSecret(db, encrypt)

	// Подготовка тестовых данных
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey1 := uuid.New()
	createdAt := time.Now()
	lastUsed := time.Now().Add(time.Hour)
	secret := &model.Secret{
		ApiKey:    apiKey1,
		ApiSecret: uuid.New(),
		ClientID:  clientID,
		IsActive:  true,
		CreatedAt: createdAt,
		LastUsed:  lastUsed,
	}
	err = secretRepo.Add(ctx, secret)
	require.NoError(t, err)

	apiKey2 := uuid.New()
	Secret2 := &model.Secret{
		ApiKey:    apiKey2,
		ApiSecret: uuid.New(),
		ClientID:  clientID,
		IsActive:  true,
	}
	err = secretRepo.Add(ctx, Secret2)
	require.NoError(t, err)

	// Создаем фильтр для поиска активных секретов
	filter := &model.SecretFilter{ApiKey: &apiKey1}

	// Выполняем тестируемую функцию find
	secrets, err := secretRepo.Find(ctx, filter)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	require.Equal(t, apiKey1, secrets[0].ApiKey)
	require.Equal(t, secret.ApiSecret, secrets[0].ApiSecret)
	require.Equal(t, true, secrets[0].IsActive)
	require.WithinDuration(t, secret.CreatedAt, createdAt, time.Second)
	require.WithinDuration(t, secret.LastUsed, lastUsed, time.Second)
}

func TestSecret_FindOneNotFound(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	secretRepo := NewSecret(db, encrypt)

	apiKey := uuid.New()

	// Создаем фильтр для поиска
	filter := &model.SecretFilter{ApiKey: &apiKey}

	// Выполняем тестируемую функцию FindOne
	_, err := secretRepo.FindOne(ctx, filter)
	require.ErrorIs(t, secretNotFoundError, err)
}

func TestSecret_Update(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	secretRepo := NewSecret(db, encrypt)

	// Подготовка тестовых данных
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	createdAt := time.Now()
	lastUsed := time.Now().Add(-time.Hour)
	secret := &model.Secret{
		ApiKey:    apiKey,
		ApiSecret: uuid.New(),
		ClientID:  clientID,
		IsActive:  true,
		CreatedAt: createdAt,
		LastUsed:  lastUsed,
	}

	err = secretRepo.Add(ctx, secret)
	require.NoError(t, err)

	// Создаем фильтр и обновляемые данные
	isActive := false
	lastUsed = time.Now()
	filter := &model.SecretFilter{ApiKey: &apiKey}
	partial := &model.SecretPartial{
		IsActive: &isActive,
		LastUsed: &lastUsed,
	}

	// Выполняем тестируемую функцию Update
	err = secretRepo.Update(ctx, secret, filter, partial)
	require.NoError(t, err)

	// Проверка обновления данных
	one, err := secretRepo.FindOne(ctx, filter)
	require.NoError(t, err)
	require.Equal(t, secret.ApiKey, one.ApiKey)
	require.Equal(t, secret.ApiSecret, one.ApiSecret)
	require.Equal(t, secret.ClientID, one.ClientID)
	require.Equal(t, isActive, one.IsActive)
	require.WithinDuration(t, secret.CreatedAt, createdAt, time.Second)
	require.WithinDuration(t, secret.LastUsed, lastUsed, time.Second)
}
