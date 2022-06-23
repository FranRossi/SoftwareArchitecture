package repository

import (
	"context"
	m "electoral_service/models"
	"encrypt"
	"fmt"
	"log"
	"strconv"
	"voter_api/connections"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindVoter(idVoter string) (*m.VoterModel, error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no voter allowed to vote with that id")
		if err2 == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err2)
	}
	other := result["otherFields"].(bson.M)
	user := &m.VoterModel{
		Id:                   result["id"].(string),
		FullName:             result["name"].(string),
		Sex:                  result["sex"].(string),
		BirthDate:            result["birthDate"].(string),
		Phone:                result["phone"].(string),
		Email:                result["email"].(string),
		Voted:                int(result["voted"].(int32)),
		LastCandidateVotedId: result["lastCandidateVotedId"].(string),
		Region:               result["region"].(string),
		OtherFields:          other,
	}
	encrypt.DecryptVoter(user)
	return user, nil
}

func FindCandidate(idCandidate string) (string, error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_votes")
	uruguayCollection := votesDatabase.Collection("votes_per_candidate")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idCandidate}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no candidate with such id")
		if err2 == mongo.ErrNoDocuments {
			return "", err2
		}
		log.Fatal(err2)
	}
	return result["id"].(string), nil
}

func FindElectionMode(idElection string) (string, error) {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayCollection := electionDatabase.Collection("configuration_election")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idElection}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("wrong election mode")
		if err2 == mongo.ErrNoDocuments {
			return "", err2
		}
		log.Fatal(err2)
	}
	return result["electionMode"].(string), nil
}

func FindElectionTime(idElection string) (string, string, error) {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayCollection := electionDatabase.Collection("configuration_election")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idElection}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("wrong election mode")
		if err2 == mongo.ErrNoDocuments {
			return "", "", err2
		}
		log.Fatal(err2)
	}
	return result["startingDate"].(string), result["finishingDate"].(string), nil
}

func HowManyVotesHasAVoter(idVoter string) int {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no voter habilitated to vote with that id")
		if err2 == mongo.ErrNoDocuments {
			return 0
		}
	}
	return int(result["voted"].(int32))
}

func GetMaxVotesAndEmailsBeforeAlert(idElection string) (int, []string, error) {
	client := connections.GetInstanceMongoClient()
	electionDatabase := client.Database("uruguay_election")
	uruguayCollection := electionDatabase.Collection("configuration_election")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idElection}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("wrong election mode")
		if err2 == mongo.ErrNoDocuments {
			return 0, []string{}, err2
		}
		log.Fatal(err2)
	}
	maxVotesString := result["otherField"].(bson.M)["maxVotes"].(string)
	maxVotes, err3 := strconv.Atoi(maxVotesString)
	emails := result["otherField"].(bson.M)["emails"].(bson.A)
	var emailsArray []string
	for _, email := range emails {
		emailsArray = append(emailsArray, email.(string))
	}
	if err3 != nil {
		return -1, []string{}, fmt.Errorf("wrong maximum values")
	}
	return maxVotes, emailsArray, nil
}
