package main

import (
	"context"
	"fmt"
	"log"
	"time"

	api_voter "239850_221025_219401/api_voter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	api_voter.Connection()
	//mongoDB()
}

func mongoDB() {
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
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	addExample(client, ctx)
	fmt.Println(databases, "anduvo")
}

func addExample(client *mongo.Client, ctx context.Context) {
	votesDatabase := client.Database("votes")
	uruguayCollection := votesDatabase.Collection("uruguayVotes")
	voteResult, err := uruguayCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"id", "1234567-8"},
		},
		bson.D{
			{"id", "1234567-9"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %v documents into uruguayVotes collection!\n", len(voteResult.InsertedIDs))

}
