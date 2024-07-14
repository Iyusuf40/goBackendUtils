package storage

import (
	"context"
	"time"

	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/redis/go-redis/v9"
)

type RedisWrapper struct {
	client *redis.Client
}

func (RW *RedisWrapper) New(db int) TempStore {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		DB:       db,
	})

	RW.client = client
	return RW
}

func (RW *RedisWrapper) GetVal(key string) string {
	return RW.client.Get(context.Background(), key).Val()
}

func (RW *RedisWrapper) SetKeyToVal(key string, value string) bool {
	return RW.client.Set(context.Background(), key, value, 0).Err() == nil
}

// expiry is in seconds
func (RW *RedisWrapper) setKeyToExpiry(key string, expiry float64) bool {
	return RW.client.Expire(context.Background(), key, time.Duration(expiry)*time.Second).Err() == nil
}

func (RW *RedisWrapper) SetKeyToValWIthExpiry(key string, value string, expiry float64) bool {
	return RW.client.SetEx(context.Background(), key, value, time.Duration(expiry)*time.Second).Err() == nil
}

func (RW *RedisWrapper) ChangeKeyEpiry(key string, newExpiry float64) bool {
	RW.setKeyToExpiry(key, newExpiry)
	return true
}

func (RW *RedisWrapper) DelKey(key string) bool {
	return RW.client.Del(context.Background(), key).Err() == nil
}

func (RW *RedisWrapper) FlushDB() {
	err := RW.client.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

func MakeRedisWrapper(db int) TempStore {
	return new(RedisWrapper).New(db)
}
