package cache

import (
	"context"
	"encoding/json"
	l "own_logger"
	"time"
)

const DefaultExpiration time.Duration = time.Second * 30

func Get(key string, result any) error {
	cache, err := GetInstanceRedisClient()
	if err != nil {
		return err
	}
	ctx := context.Background()

	valInBytes, err := cache.Get(ctx, key).Bytes()
	if err != nil {
		go l.LogInfo(" Cache: Value with key  " + key + " was not found in cache")
		return err
	}

	go l.LogInfo("Cache: Value with key " + key + " was found in cache")
	err = json.Unmarshal(valInBytes, &result)
	return err
}

func Save(key string, value any, expiration time.Duration) error {

	go l.LogInfo("Saving value with key  " + key + " to cache")
	cache, err := GetInstanceRedisClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	valInBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, errSet := cache.Set(ctx, key, valInBytes, expiration).Result()

	return errSet
}
