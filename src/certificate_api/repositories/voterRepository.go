package repositories

import (
	"certificate_api/connections"
	"context"
	electoral_service_models "electoral_service/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (certRepo *CertificatesRepo) FindVoter(idVoter string) (user electoral_service_models.VoterModel, err error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		return user, fmt.Errorf("there is no voter assigned to vote with that id: %v", err2)
	}
	other := result["otherFields"].(bson.M)
	user = electoral_service_models.VoterModel{
		Id:          result["id"].(string),
		FullName:    result["name"].(string),
		Sex:         result["sex"].(string),
		BirthDate:   result["birthDate"].(string),
		Phone:       result["phone"].(string),
		Email:       result["email"].(string),
		Voted:       int(result["voted"].(int32)),
		Region:      result["region"].(string),
		OtherFields: other,
	}
	return user, nil
}

func (certRepo *CertificatesRepo) FindElection(idElection string) (election electoral_service_models.ElectionModelEssential, err error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("configuration_election")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idElection}}).Decode(&result)
	if err2 != nil {
		return election, fmt.Errorf("there is no election with that id: %v", err2)
	}
	election = electoral_service_models.ElectionModelEssential{
		Id:            result["id"].(string),
		ElectionMode:  result["electionMode"].(string),
		StartingDate:  result["startingDate"].(string),
		FinishingDate: result["configuration_election"].(string),
	}
	return election, nil
}
