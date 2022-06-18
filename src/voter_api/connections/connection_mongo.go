package connections

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	l "own_logger"
	"sync"
	"time"
)

func connectionMongo() *mongo.Client {
	address := os.Getenv("mongo_address")
	client, err := mongo.NewClient(options.Client().ApplyURI(address))

	if err != nil {
		l.LogError(err.Error())
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		l.LogError(err.Error())
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		l.LogError(err.Error())
	}
	return client
}

var lock = &sync.Mutex{}

var mongoClientInstance *mongo.Client

func GetInstanceMongoClient() *mongo.Client {
	if mongoClientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if mongoClientInstance == nil {
			fmt.Println("Creating mongo client instance now.")
			mongoClientInstance = connectionMongo()
		}
	}
	return mongoClientInstance
}
