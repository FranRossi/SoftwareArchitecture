package repositories

import (
	"certificate_api/connections"
	"context"
	"fmt"
	"log"

	electoral_service_models "electoral_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindVoterFullName(idVoter string) (fullName string, err error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no voter habilitated to vote with that id")
		if err2 == mongo.ErrNoDocuments {
			return "", nil
		}
		log.Fatal(err2)
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
		OtherFields: other,
	}
	return user.FullName, nil
}
