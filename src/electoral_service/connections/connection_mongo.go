package connections

import (
	"context"
	"fmt"
	"log"
	l "own_logger"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lockDB = &sync.Mutex{}

var mongoClientInstance *mongo.Client

func GetInstanceMongoClient() *mongo.Client {
	if mongoClientInstance == nil {
		lockDB.Lock()
		defer lockDB.Unlock()
		if mongoClientInstance == nil {
			fmt.Println("Creating mongo client instance now.")
			mongoClientInstance = connectionMongo()
		}
	}
	return mongoClientInstance
}

func connectionMongo() *mongo.Client {
	const uri = "mongodb://localhost:27017" // TODO .env
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client().ApplyURI(uri).SetDirect(true)
	client, err := mongo.Connect(ctx, opts)
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
