package pg

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClient_Add(t *testing.T) {
	ctx := context.Background()

	defer func(ctx context.Context) {
		err := truncateAll(ctx)
		require.NoError(t, err)
	}(ctx)

	// Создаем клиента для работы с базой данных
	clientRepo := NewClient(db)

	// Пример данных клиента для добавления
	client := &model.Client{
		ID:   uuid.New(),
		Name: "Test Client",
	}

	// Выполняем тестируемую функцию Add
	err := clientRepo.Add(ctx, client)
	require.NoError(t, err)

	// Проверяем, что данные успешно добавлены
	var id uuid.UUID
	var name string
	err = db.QueryRow(ctx, `SELECT id, name FROM endpoint.clients WHERE id=$1`, client.ID).Scan(&id, &name)
	require.NoError(t, err)
	require.Equal(t, client.ID, id)
	require.Equal(t, client.Name, name)
}
