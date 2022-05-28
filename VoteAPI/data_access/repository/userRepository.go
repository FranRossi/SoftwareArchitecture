package repository

import (
	domain "VoteAPI/domain/user"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func RegisterUser(user *domain.User) error {
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
	usersDatabase := client.Database("users")
	uruguayVotersCollection := usersDatabase.Collection("uruguayVoters")
	_, err2 := uruguayVotersCollection.InsertOne(ctx, bson.M{"id": user.Id, "username": user.Username, "password": user.HashedPassword, "token": user.Token})
	if err2 != nil {
		fmt.Println("error creating user")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return err2
}
