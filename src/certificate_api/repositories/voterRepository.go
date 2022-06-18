package repositories

import (
	"certificate_api/connections"
	"context"
	electoral_service_models "electoral_service/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func FindVoterFullName(idVoter string) (fullName string, err error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		return "", fmt.Errorf("there is no voter assigned to vote with that id: %v", err2)
	}
	other := result["otherFields"].(bson.M)
	user := &electoral_service_models.VoterModel{
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
	return user.FullName, nil
}
