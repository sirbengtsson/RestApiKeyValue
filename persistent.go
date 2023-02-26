package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type PersistentStore struct {
	redisClient *redis.Client
}

var ctx = context.Background()

func (store *PersistentStore) setValue(key string, value string) error {
	err := store.redisClient.Set(ctx, key, value, 0).Err()

	return err
}

func (store *PersistentStore) readValue(key string) (string, error) {
	value, err := store.redisClient.Get(ctx, key).Result()

	return value, err
}

func (store *PersistentStore) deleteValue(key string) error {
	_, err := store.redisClient.Del(ctx, key).Result()

	return err
}

func (store *PersistentStore) getAllKeys() ([]string, error) {
	dbSize, err := store.redisClient.DBSize(ctx).Result()

	if err != nil {
		return []string{}, err
	}

	keys := make([]string, 0, dbSize)

	iter := store.redisClient.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	return keys, iter.Err()
}
