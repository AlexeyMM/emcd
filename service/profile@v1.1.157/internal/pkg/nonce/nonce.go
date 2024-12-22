package nonce

import (
	"context"
	"sync"
)

type Store struct {
	sync.RWMutex
	data map[string]int64
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]int64),
	}
}

func (store *Store) CheckAndUpdateNonce(ctx context.Context, userID string, newNonce int64) (bool, error) {
	store.Lock()
	defer store.Unlock()

	lastNonce, exists := store.data[userID]

	if exists && newNonce <= lastNonce {
		return false, nil
	}

	store.data[userID] = newNonce
	return true, nil
}

//TODO: заменить на Redis
//type RedisStore struct {
//	client *redis.Client
//}
//
//func NewRedisStore(addr string, password string, db int) *RedisStore {
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     addr,
//		Password: password,
//		DB:       db,
//	})
//
//	return &RedisStore{
//		client: rdb,
//	}
//}
//
//func (store *RedisStore) CheckAndUpdateNonce(ctx context.Context, userID, apiKey string, newNonce int64) (bool, error) {
//	key := userID
//
//	script := redis.NewScript(`
//		local lastNonce = redis.call("GET", KEYS[1])
//		if lastNonce then
//			if tonumber(lastNonce) >= tonumber(ARGV[1]) then
//				return 0
//			end
//		end
//		redis.call("SET", KEYS[1], ARGV[1])
//		return 1
//	`)
//
//	result, err := script.Run(ctx, store.client, []string{key}, newNonce).Result()
//	if err != nil {
//		return false, fmt.Errorf("redis check and update nonce: %w", err)
//	}
//
//	if result.(int64) == 1 {
//		return true, nil
//	}
//
//	return false, nil
//}
