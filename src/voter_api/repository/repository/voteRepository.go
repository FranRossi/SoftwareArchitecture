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
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes")
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
	totalVotesCollection := electionDatabase.Collection("total_votes")
	filter = bson.D{{"id", vote.IdElection}}
	update = bson.D{{"$inc", bson.D{{"votes_counted", 1}}}}
	_, err2 = totalVotesCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error registering new vote on election")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}

func RegisterVote(idVoter string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	filter := bson.D{{"id", idVoter}}
	update := bson.D{{"$inc", bson.D{{"voted", 1}}}}
	_, err2 := uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error registering new vote for candidate")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
