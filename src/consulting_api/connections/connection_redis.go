package connections

import (
	"context"
	"fmt"
	l "own_logger"
	"sync"

	"github.com/go-redis/redis/v8"
)

var lockRedis = &sync.Mutex{}

var redisClientInstance *redis.Client

func GetInstanceRedisClient() *redis.Client {
	if redisClientInstance == nil {
		lockRedis.Lock()
		defer lock.Unlock()
		if redisClientInstance == nil {
			fmt.Println("Creating redis client instance now.")
			redisClientInstance, _ = connectionRedis()
		}
	}
	return redisClientInstance
}

func connectionRedis() (*redis.Client, error) {
	const uri = "localhost:6379"
	const password = "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81" // TODO: use .env file

	client := redis.NewClient(&redis.Options{Addr: uri, Password: password, DB: 0})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		l.LogError(err.Error())
		return nil, err
	}
	return client, nil
}
