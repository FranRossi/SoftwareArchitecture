package repository

import (
	"context"
	models2 "electoral_api/adapter/uruguayan_election/models"
	"electoral_api/connections"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ElectionRepo struct {
}

func (repo *ElectionRepo) StoreElectionConfiguration(election models2.ElectionModel) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguayan_election")
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
	uruguayDataBase := client.Database("uruguayan_election")
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
