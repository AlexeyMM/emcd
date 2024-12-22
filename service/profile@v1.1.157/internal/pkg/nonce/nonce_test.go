package nonce

import (
	"context"
	"testing"
)

func TestStore_CheckAndUpdateNonce(t *testing.T) {
	store := NewStore()
	ctx := context.Background()

	// Тест 1: Первый nonce для пользователя
	t.Run("First nonce", func(t *testing.T) {
		userID := "user1"
		newNonce := int64(100)

		success, err := store.CheckAndUpdateNonce(ctx, userID, newNonce)
		if err != nil {
			t.Fatalf("Ошибка при обновлении nonce: %v", err)
		}

		if !success {
			t.Errorf("Ожидалось, что обновление пройдет успешно, но оно не прошло")
		}
	})

	// Тест 2: Новый nonce больше предыдущего
	t.Run("Nonce is increasing", func(t *testing.T) {
		userID := "user1"
		newNonce := int64(200)

		success, err := store.CheckAndUpdateNonce(ctx, userID, newNonce)
		if err != nil {
			t.Fatalf("Ошибка при обновлении nonce: %v", err)
		}

		if !success {
			t.Errorf("Ожидалось, что обновление пройдет успешно, но оно не прошло")
		}
	})

	// Тест 3: Новый nonce меньше предыдущего
	t.Run("Nonce is not increasing", func(t *testing.T) {
		userID := "user1"
		newNonce := int64(150)

		success, err := store.CheckAndUpdateNonce(ctx, userID, newNonce)
		if err != nil {
			t.Fatalf("Ошибка при обновлении nonce: %v", err)
		}

		if success {
			t.Errorf("Ожидалось, что обновление не пройдет, но оно прошло")
		}
	})

	// Тест 4: Новый nonce равен предыдущему
	t.Run("Nonce is equal to previous", func(t *testing.T) {
		userID := "user1"
		newNonce := int64(200)

		success, err := store.CheckAndUpdateNonce(ctx, userID, newNonce)
		if err != nil {
			t.Fatalf("Ошибка при обновлении nonce: %v", err)
		}

		if success {
			t.Errorf("Ожидалось, что обновление не пройдет, но оно прошло")
		}
	})

	// Тест 5: Проверка с другим пользователем
	t.Run("Different user", func(t *testing.T) {
		userID := "user2"
		newNonce := int64(50)

		success, err := store.CheckAndUpdateNonce(ctx, userID, newNonce)
		if err != nil {
			t.Fatalf("Ошибка при обновлении nonce: %v", err)
		}

		if !success {
			t.Errorf("Ожидалось, что обновление пройдет успешно, но оно не прошло")
		}
	})
}
