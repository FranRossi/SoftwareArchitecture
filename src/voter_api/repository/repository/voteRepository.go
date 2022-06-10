package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"voter_api/connections"
	domain "voter_api/domain/vote"
)

func StoreVote(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	usersDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := usersDatabase.Collection("votes")
	filter := bson.D{{"id", vote.IdCandidate}}
	update := bson.D{{"$inc", bson.D{{"votes", 1}}}}
	_, err2 := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error storing vote")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
