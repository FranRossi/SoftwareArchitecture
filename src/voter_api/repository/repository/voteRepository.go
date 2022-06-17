package repository

import (
	"context"
	"encrypt"
	"fmt"
	"log"
	"voter_api/connections"
	"voter_api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StoreVote(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes_per_candidate")
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
	filter2 := bson.D{{"election_id", vote.IdElection}}
	update2 := bson.D{{"$inc", bson.D{{"votes_counted", 1}}}}
	_, err3 := totalVotesCollection.UpdateOne(context.TODO(), filter2, update2)
	if err3 != nil {
		fmt.Println("error registering new vote on election")
		if err3 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err3)
	}
	return nil
}

func RegisterVote(vote *domain.VoteModel, electionMode string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$inc", bson.D{{"voted", 1}}}}
	_, err2 := uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		message := "error registering new vote for candidate"
		if err2 == mongo.ErrNoDocuments {
			return fmt.Errorf("voter not found")
		}
		log.Fatal(err2)
		return fmt.Errorf(message)
	}
	if electionMode == "multi" {
		err := setCandidateToVoter(vote)
		if err != nil {
			return err
		}
	}
	return nil
}

func setCandidateToVoter(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$set", bson.D{{"lastCandidateVotedId", encrypt.EncryptText(vote.IdCandidate)}}}}
	_, err2 := uruguayCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err2 != nil {
		message := "error registering last candidate for voter"
		if err2 == mongo.ErrNoDocuments {
			return fmt.Errorf("not documents found")
		}
		log.Fatal(err2)
		return fmt.Errorf(message)
	}
	return nil
}

func DeleteVote(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes_per_candidate")
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
	filter = bson.D{{"election_id", vote.IdElection}}
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

func StoreVoteInfo(idVoter, idElection, timeFrontEnd, timeBackEnd, timePassed, voteIdentification string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayCollection := uruguayDataBase.Collection("votes_info")
	_, err := uruguayCollection.InsertOne(context.TODO(), bson.D{{"voter", idVoter}, {"election", idElection}, {"time_front_end", timeFrontEnd}, {"time_back_end", timeBackEnd}, {"time_passed", timePassed}, {"vote_identification", voteIdentification}})
	if err != nil {
		return fmt.Errorf("error storing vote info: %v", err)
	}
	return nil
}

func ReplaceLastCandidateVoted(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayVotersCollection := electionDatabase.Collection("voters")
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$set", bson.D{{"lastCandidateVotedId", encrypt.EncryptText(vote.IdCandidate)}}}}
	_, err2 := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		message := "error replacing candidate for voter"
		if err2 == mongo.ErrNoDocuments {
			return fmt.Errorf("not documents found")
		}
		log.Fatal(err2)
		return fmt.Errorf(message)
	}
	return nil
}

func DeleteOldVote(vote *domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayVotersCollection := electionDatabase.Collection("voters")
	var result bson.M
	err := uruguayVotersCollection.FindOne(context.TODO(), bson.D{{"id", vote.IdVoter}}).Decode(&result)
	if err != nil {
		return fmt.Errorf("error deleting old vote: %v", err)
	}
	lastCandidateVotedId := encrypt.DecryptText(result["lastCandidateVotedId"].(string))

	votesDatabase := client.Database("uruguay_votes")
	votesPerCandidatesCollection := votesDatabase.Collection("votes_per_candidate")
	filter2 := bson.D{{"id", lastCandidateVotedId}}
	update2 := bson.D{{"$inc", bson.D{{"votes", -1}}}}
	_, err2 := votesPerCandidatesCollection.UpdateOne(context.TODO(), filter2, update2)
	if err2 != nil {
		message := "error deleting old vote from candidate"
		if err2 == mongo.ErrNoDocuments {
			return fmt.Errorf("not documents found")
		}
		log.Fatal(err2)
		return fmt.Errorf(message)
	}
	return nil
}
