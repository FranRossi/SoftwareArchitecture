package data_access

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoConnectionParameters struct {
	Client *mongo.Client
	Ctx    context.Context
}

var mongoParameters *MongoConnectionParameters

func MongoConnection() {
	const uri = "mongodb://localhost:27017"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	*mongoParameters = MongoConnectionParameters{
		Client: client,
		Ctx:    ctx,
	}
}

func GetConnectionParameters() *MongoConnectionParameters {
	return mongoParameters
}
