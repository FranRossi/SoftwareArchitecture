package cache

import (
	"context"
	"os"
	"sync"

	l "own_logger"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var lockRedis = &sync.Mutex{}

var redisClientInstance *redis.Client

func GetInstanceRedisClient() (*redis.Client, error) {

	var err error
	if redisClientInstance == nil {
		lockRedis.Lock()
		defer lockRedis.Unlock()
		if redisClientInstance == nil {
			l.LogInfo("Creating redis client instance now.")
			redisClientInstance, err = connectionRedis()
		}
	}
	return redisClientInstance, err
}

func connectionRedis() (*redis.Client, error) {
	errLoadingEnv := godotenv.Load("./../cache/.env")
	if errLoadingEnv != nil {
		l.LogError("Error loading .env file of cache: " + errLoadingEnv.Error())
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
