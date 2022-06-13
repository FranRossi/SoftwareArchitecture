package connections

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"sync"
	"time"
)

func connectionMongo() *mongo.Client {
	address := os.Getenv("mongo_address")
	client, err := mongo.NewClient(options.Client().ApplyURI(address))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
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
