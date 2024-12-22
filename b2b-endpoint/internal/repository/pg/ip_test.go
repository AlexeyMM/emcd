package pg

import (
	"context"
	"testing"
	"time"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIP_Add(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	// Пример данных для клиента и секретов
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO endpoint.secrets (api_key, api_secret, client_id) VALUES ($1, $2, $3)`, apiKey, "test_secret", clientID)
	require.NoError(t, err)

	ipRepo := NewIP(db)

	// Пример данных IP для добавления
	ip1 := &model.IP{
		ApiKey:    apiKey,
		Address:   "192.168.1.1",
		CreatedAt: time.Now().UTC(),
	}
	ip2 := &model.IP{
		ApiKey:    apiKey,
		Address:   "192.168.1.2",
		CreatedAt: time.Now().UTC(),
	}

	ips := []*model.IP{ip1, ip2}

	err = ipRepo.Add(ctx, ips)
	require.NoError(t, err)

	// Проверяем, что данные успешно добавлены
	found, err := ipRepo.Find(ctx, &model.IPFilter{
		ApiKey: &apiKey,
	})
	require.NoError(t, err)
	require.Len(t, found, len(ips))

	for _, ip := range found {
		switch ip.Address {
		case ip1.Address:
			require.Equal(t, ip.ApiKey, ip1.ApiKey)
			require.Equal(t, ip.Address, ip1.Address)
			require.WithinDuration(t, ip.CreatedAt, ip1.CreatedAt, time.Second)
		case ip2.Address:
			require.Equal(t, ip.ApiKey, ip2.ApiKey)
			require.Equal(t, ip.Address, ip2.Address)
			require.WithinDuration(t, ip.CreatedAt, ip2.CreatedAt, time.Second)
		}
	}
}

func TestIP_Delete(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	ipRepo := NewIP(db)

	// Пример данных для клиента и секретов
	clientID := uuid.New()
	_, err := db.Exec(ctx, `INSERT INTO endpoint.clients (id, name) VALUES ($1, $2)`, clientID, "Test Client")
	require.NoError(t, err)

	apiKey := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO endpoint.secrets (api_key, api_secret, client_id) VALUES ($1, $2, $3)`, apiKey, "test_secret", clientID)
	require.NoError(t, err)

	// Добавляем IP для удаления
	ip := &model.IP{
		ApiKey:    apiKey,
		Address:   "192.168.1.1",
		CreatedAt: time.Now(),
	}

	ips := []*model.IP{ip}
	err = ipRepo.Add(ctx, ips)
	require.NoError(t, err)

	// Удаляем IP
	filter := &model.IPFilter{ApiKey: &apiKey, Address: &ip.Address}
	err = ipRepo.Delete(ctx, filter)
	require.NoError(t, err)

	// Проверяем, что IP-адрес удалён
	var count int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM endpoint.whitelist_ips WHERE api_key=$1`, ip.ApiKey).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}
