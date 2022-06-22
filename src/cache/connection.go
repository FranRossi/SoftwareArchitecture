package cache

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	l "own_logger"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var lockRedis = &sync.Mutex{}

var redisClientInstance *redis.Client

var cacheDisable = false

var timeForDisablingCache time.Time

func GetInstanceRedisClient() (*redis.Client, error) {
	if cacheDisable {
		if time.Now().After(timeForDisablingCache) {
			cacheDisable = false
			l.LogInfo("Enabling cache")
		} else {
			return nil, errors.New("cache was disable by the programme")
		}
	}
	var err error
	if redisClientInstance == nil {
		lockRedis.Lock()
		defer lockRedis.Unlock()
		if redisClientInstance == nil {
			l.LogInfo("Creating redis client instance now.")
			redisClientInstance, err = connectionRedis()
		}
	}
	if err != nil {
		cacheDisable = true
		timeForDisablingCache = time.Now().Add(1 * time.Minute)
		l.LogWarning("Disabling  cache, connection failed")
	}
	return redisClientInstance, err
}

func connectionRedis() (*redis.Client, error) {
	errLoadingEnv := godotenv.Load("./../cache/.env")
	if errLoadingEnv != nil {
		l.LogError("Error loading .env file of cache: " + errLoadingEnv.Error())
		return nil, errLoadingEnv
	}
	uri := os.Getenv("REDIS_URI")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{Addr: uri, Password: password, DB: 0})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		l.LogError(err.Error())
		return nil, err
	}
	return client, nil
}
