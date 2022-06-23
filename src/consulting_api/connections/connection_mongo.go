package connections

import (
	"context"
	"fmt"
	"log"
	"os"
	l "own_logger"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

func connectionMongo() *mongo.Client {
	uri := os.Getenv("MONGO")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

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
		log.Fatal(err)
	}
	return client
}

func CloseMongoClient() {
	if mongoClientInstance != nil {
		err := mongoClientInstance.Disconnect(context.TODO())
		if err != nil {
			l.LogError(err.Error() + " cannot close mongo client")
			return
		}
		l.LogInfo("mongo client closed")
	}
}
