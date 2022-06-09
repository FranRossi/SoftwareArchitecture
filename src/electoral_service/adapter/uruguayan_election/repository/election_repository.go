package repository

import (
	"context"
	models2 "electoral_service/adapter/uruguayan_election/models"
	"electoral_service/connections"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ElectionRepo struct {
}

func (repo *ElectionRepo) StoreElectionConfiguration(election models2.ElectionModel) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayanElectionCollection := uruguayDataBase.Collection("configuration_election")
	_, err := uruguayanElectionCollection.InsertOne(context.TODO(), bson.M{"id": election.Id, "description": election.Description, "startingDate": election.StartingDate, "finishingDate": election.FinishingDate, "electionMode": election.ElectionMode})
	if err != nil {
		fmt.Println("error storing election configuration")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return err
}

func StoreElectionVoters(voters []models2.VoterModel) error {
	votersInterface := convertModelToInterface(voters)
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayanVotersCollection := uruguayDataBase.Collection("voters")
	_, err := uruguayanVotersCollection.InsertMany(context.TODO(), votersInterface)
	if err != nil {
		fmt.Println("error storing voters")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return err
}

func convertModelToInterface(voters []models2.VoterModel) []interface{} {
	var votersInterface []interface{}

	for _, v := range voters {
		votersInterface = append(votersInterface, v)
	}
	return votersInterface
}

func StoreCandidates(candidates []models2.CandidateModel) error {
	client := connections.GetInstanceMongoClient()
	candidatesToStore := convertCandidateModelToInterface(candidates)
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanCandidatesCollection := uruguayDataBase.Collection("votes")
	_, err := uruguayanCandidatesCollection.InsertMany(context.TODO(), candidatesToStore)
	if err != nil {
		fmt.Println("error storing candidates")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return err
}

type Candidate struct {
	Id    string `bson:"id"`
	Name  string `bson:"name"`
	Votes int    `bson:"votes"`
}

func convertCandidateModelToInterface(candidates []models2.CandidateModel) []interface{} {
	var candidatesResume []Candidate
	for _, v := range candidates {
		candidatesResume = append(candidatesResume, Candidate{Id: v.Id, Name: v.Name, Votes: 0})
	}

	var candidatesInterface []interface{}
	for _, v := range candidatesResume {
		candidatesInterface = append(candidatesInterface, v)
	}
	return candidatesInterface
}
