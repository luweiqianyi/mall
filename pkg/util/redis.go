package util

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mall/pkg/log"
	"time"
)

func NewRedisClient(RedisHost, RedisPassword string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: RedisPassword,
		DB:       1,
	})
	return rdb
}

func RedisSet(rdb *redis.Client, key string, value interface{}, expiration time.Duration) error {
	if rdb == nil {
		return fmt.Errorf("redis client nil")
	}
	err := rdb.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.PrintLog("%v", err)
		return err
	}
	return nil
}

func RedisGet(rdb *redis.Client, key string) (string, error) {
	if rdb == nil {
		return "", fmt.Errorf("redis client nil")
	}
	value, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.PrintLog("%v", err)
		return "", err
	}
	return value, nil
}

func RedisDel(rdb *redis.Client, key string) error {
	if rdb == nil {
		return fmt.Errorf("redis client nil")
	}
	err := rdb.Del(context.Background(), key).Err()
	if err != nil {
		log.PrintLog("%v", err)
		return err
	}
	return nil
}
