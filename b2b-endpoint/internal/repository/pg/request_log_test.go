package pg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRequestLog_Add(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	logRepo := NewRequestLog(db)

	// Пример данных для клиента и секретов
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO endpoint.secrets (api_key, api_secret, client_id) VALUES ($1, $2, $3)`, apiKey, "test_secret", clientID)
	require.NoError(t, err)

	// Пример данных для лога
	requestLog := &model.RequestLog{
		ApiKey:      apiKey,
		RequestHash: "test_request_hash",
		CreatedAt:   time.Now(),
	}

	// Выполняем тестируемую функцию Add
	err = logRepo.Add(ctx, requestLog)
	require.NoError(t, err)

	// Проверяем, что данные успешно добавлены
	var apiKeyResult uuid.UUID
	var requestHash string
	var createdAt time.Time
	err = db.QueryRow(ctx, `SELECT api_key, request_hash, created_at FROM endpoint.request_logs WHERE api_key=$1`, requestLog.ApiKey).
		Scan(&apiKeyResult, &requestHash, &createdAt)
	require.NoError(t, err)
	require.Equal(t, requestLog.ApiKey, apiKeyResult)
	require.Equal(t, requestLog.RequestHash, requestHash)
	require.WithinDuration(t, requestLog.CreatedAt, createdAt, time.Second)
}

func TestRequestLog_FindOne(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	logRepo := NewRequestLog(db)

	// Пример данных для клиента и секретов
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO endpoint.secrets (api_key, api_secret, client_id) VALUES ($1, $2, $3)`, apiKey, "test_secret", clientID)
	require.NoError(t, err)

	// Пример данных для лога
	hash := "test_request_hash"
	requestLog := &model.RequestLog{
		ApiKey:      apiKey,
		RequestHash: hash,
		CreatedAt:   time.Now(),
	}
	err = logRepo.Add(ctx, requestLog)
	require.NoError(t, err)

	// Выполняем тестируемую функцию FindOne
	filter := &model.RequestLogFilter{ApiKey: &apiKey, RequestHash: &hash}
	foundLog, err := logRepo.FindOne(ctx, filter)
	require.NoError(t, err)
	require.Equal(t, requestLog.ApiKey, foundLog.ApiKey)
	require.Equal(t, requestLog.RequestHash, foundLog.RequestHash)
	require.WithinDuration(t, requestLog.CreatedAt, foundLog.CreatedAt, time.Second)
}

func TestRequestLog_FindOneNotFound(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	logRepo := NewRequestLog(db)

	apiKey := uuid.New()
	hash := "test_request_hash"
	// Выполняем тестируемую функцию FindOne
	filter := &model.RequestLogFilter{ApiKey: &apiKey, RequestHash: &hash}
	_, err := logRepo.FindOne(ctx, filter)
	require.ErrorIs(t, LogsNotFoundError, err)
}

func TestRequestLog_Find(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	logRepo := NewRequestLog(db)

	// Пример данных для клиента и секретов
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO endpoint.secrets (api_key, api_secret, client_id) VALUES ($1, $2, $3)`, apiKey, "test_secret", clientID)
	require.NoError(t, err)

	// Пример данных для логов
	for i := 1; i <= 3; i++ {
		requestLog := &model.RequestLog{
			ApiKey:      apiKey,
			RequestHash: fmt.Sprintf("test_request_hash_%d", i),
			CreatedAt:   time.Now(),
		}
		err = logRepo.Add(ctx, requestLog)
		require.NoError(t, err)
	}

	// Выполняем тестируемую функцию Find
	filter := &model.RequestLogFilter{ApiKey: &apiKey}
	foundLogs, err := logRepo.Find(ctx, filter)
	require.NoError(t, err)
	require.Len(t, foundLogs, 3)

	for i, log := range foundLogs {
		require.Equal(t, apiKey, log.ApiKey)
		require.Equal(t, fmt.Sprintf("test_request_hash_%d", i+1), log.RequestHash)
	}
}
