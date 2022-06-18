package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"voter_api/connections"
	"voter_api/domain"
)

func StoreVote(vote domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes_per_candidate")
	filter := bson.D{{"id", vote.IdCandidate}}
	update := bson.D{{"$inc", bson.D{{"votes", 1}}}}
	_, err := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error storing vote: %v", err)
	}

	totalVotesCollection := electionDatabase.Collection("total_votes")
	filter2 := bson.D{{"election_id", vote.IdElection}}
	update2 := bson.D{{"$inc", bson.D{{"votes_counted", 1}}}}
	_, err2 := totalVotesCollection.UpdateOne(context.TODO(), filter2, update2)
	if err2 != nil {
		return fmt.Errorf("error registering new vote on election: %v", err2)
	}
	return nil
}

func RegisterVote(vote domain.VoteModel, electionMode string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$inc", bson.D{{"voted", 1}}}}
	_, err2 := uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		message := "error registering new vote for candidate"
		return fmt.Errorf(message+": %v", err2)
	}
	if electionMode == "multi" {
		err := setCandidateToVoter(vote)
		if err != nil {
			return err
		}
	}
	return nil
}

func setCandidateToVoter(vote domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$set", bson.D{{"lastCandidate", vote.IdCandidate}}}}
	_, err := uruguayCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		message := "error registering last candidate for voter"
		return fmt.Errorf(message+": %v", err)
	}
	return nil
}

func DeleteVote(vote domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_votes")
	uruguayVotersCollection := electionDatabase.Collection("votes_per_candidate")
	filter := bson.D{{"id", vote.IdCandidate}}
	update := bson.D{{"$inc", bson.D{{"votes", -1}}}}
	_, err := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error deleting vote from candidate: %v", err)
	}

	totalVotesCollection := electionDatabase.Collection("total_votes")
	filter = bson.D{{"election_id", vote.IdElection}}
	update = bson.D{{"$inc", bson.D{{"votes_counted", -1}}}}
	_, err2 := totalVotesCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		return fmt.Errorf("error deleting vote from total votes: %v", err2)
	}

	uruguayDataBase := client.Database("uruguay_election")
	uruguayCollection := uruguayDataBase.Collection("voters")
	filter = bson.D{{"id", vote.IdVoter}}
	update = bson.D{{"$inc", bson.D{{"voted", -1}}}}
	_, err3 := uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err3 != nil {
		return fmt.Errorf("error deleting vote from voter: %v", err3)
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

func ReplaceLastCandidateVoted(vote domain.VoteModel) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayVotersCollection := electionDatabase.Collection("voters")
	filter := bson.D{{"id", vote.IdVoter}}
	update := bson.D{{"$set", bson.D{{"lastCandidate", vote.IdCandidate}}}}
	_, err2 := uruguayVotersCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		message := "error replacing candidate for voter"
		return fmt.Errorf(message+": %v", err2)
	}
	return nil
}

func DeleteOldVote(vote domain.VoteModel, region string) error {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayVotersCollection := electionDatabase.Collection("voters")
	var result bson.M
	err := uruguayVotersCollection.FindOne(context.TODO(), bson.D{{"id", vote.IdVoter}}).Decode(&result)
	if err != nil {
		return fmt.Errorf("error deleting old vote: %v", err)
	}
	lastCandidateVoted := result["lastCandidate"].(string)

	votesDatabase := client.Database("uruguay_votes")
	votesPerCandidatesCollection := votesDatabase.Collection("votes_per_candidate")
	filter2 := bson.D{{"id", lastCandidateVoted}}
	update2 := bson.D{{"$inc", bson.D{{"votes", -1}}}}
	var politicalPartyObject bson.M
	err2 := votesPerCandidatesCollection.FindOneAndUpdate(context.TODO(), filter2, update2).Decode(&politicalPartyObject)
	if err2 != nil {
		message := "error deleting old vote from candidate"
		return fmt.Errorf(message+": %v", err2)
	}
	latestPoliticalParty := politicalPartyObject["politicalParty"].(string)

	uruguayanVotesCollection := votesDatabase.Collection("result_election")
	query := bson.M{"election_id": vote.IdElection}
	updateDocument := bson.M{"$inc": bson.M{
		"votes_per_candidates.$[candidate].votes": -1,
		"votes_per_parties.$[party].votes":        -1,
		"regions.$[region].votes":                 -1,
		"amount_voted":                            -1,
	}}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.D{{"candidate.id", lastCandidateVoted}}, bson.D{{"party.name", latestPoliticalParty}}, bson.D{{"region.name", region}}},
	})

	_, err = uruguayanVotesCollection.UpdateOne(context.TODO(), query, updateDocument, opts)
	if err != nil {
		return fmt.Errorf("error deleting vote from result election: %v", err)
	}
	return nil
}

func UpdateElectionResult(vote domain.VoteModel, region, politicalParty string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanVotesCollection := uruguayDataBase.Collection("result_election")

	query := bson.M{"election_id": vote.IdElection}
	updateDocument := bson.M{"$inc": bson.M{
		"votes_per_candidates.$[candidate].votes": 1,
		"votes_per_parties.$[party].votes":        1,
		"regions.$[region].votes":                 1,
		"amount_voted":                            1,
	}}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.D{{"candidate.id", vote.IdCandidate}}, bson.D{{"party.name", politicalParty}}, bson.D{{"region.name", region}}},
	})
	_, err := uruguayanVotesCollection.UpdateOne(context.TODO(), query, updateDocument, opts)

	if err != nil {
		return fmt.Errorf("error updating election result: %v", err)
	}
	return nil
}

func FindPoliticalPartyFromCandidateId(candidateId string) (string, error) {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayVotersCollection := uruguayDataBase.Collection("votes_per_candidate")
	var result bson.M
	err := uruguayVotersCollection.FindOne(context.TODO(), bson.D{{"id", candidateId}}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding political party: %v", err)
	}
	politicalParty := result["politicalParty"].(string)
	return politicalParty, nil
}
