package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func CheckVoterId(idVoter string) (bson.M, error) {
	//mongoParams := data_access.GetConnectionParameters()
	//client := mongoParams.Client
	//ctx := mongoParams.Ctx
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
	votesDatabase := client.Database("votes")
	uruguayCollection := votesDatabase.Collection("uruguayVotes")
	var result bson.M
	err2 := uruguayCollection.FindOne(ctx, bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		fmt.Println("El error esta en mongo")
		if err2 == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err2)
	}
	return result, err2
}
