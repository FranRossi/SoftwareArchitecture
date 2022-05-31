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
	usersDatabase := client.Database("votes")
	uruguayVotersCollection := usersDatabase.Collection("uruguayVotes")
	_, err2 := uruguayVotersCollection.InsertOne(context.TODO(), bson.M{"department": vote.Department, "circuit": vote.Circuit, "candidate": vote.Candidate, "politicalParty": vote.PoliticalParty})
	if err2 != nil {
		fmt.Println("error storing vote")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return err2
}
