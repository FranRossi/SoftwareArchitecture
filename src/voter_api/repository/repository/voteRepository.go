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

func DeleteVote(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes")
	filter := bson.D{{"id", vote.IdCandidate}}
	update := bson.D{{"$inc", bson.D{{"votes", -1}}}}
	_, err2 := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error deleting vote from candidate")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}

	totalVotesCollection := electionDatabase.Collection("total_votes")
	filter = bson.D{{"id", vote.IdElection}}
	update = bson.D{{"$inc", bson.D{{"votes_counted", -1}}}}
	_, err2 = totalVotesCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error deleting vote from total votes")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}

	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	filter = bson.D{{"id", vote.IdVoter}}
	update = bson.D{{"$inc", bson.D{{"voted", -1}}}}
	_, err2 = uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error deleting vote from voter ")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}

func StoreVoteInfo(idVoter, timeFrontEnd, timeBackEnd, timePassed, voteIdentification string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayCollection := uruguayDataBase.Collection("votes_info")
	_, err := uruguayCollection.InsertOne(context.TODO(), bson.D{{"voter", idVoter}, {"time_front_end", timeFrontEnd}, {"time_back_end", timeBackEnd}, {"time_passed", timePassed}, {"vote_identificator", voteIdentification}})
	if err != nil {
		return fmt.Errorf("error storing vote info: %v", err)
	}
	return nil
}
