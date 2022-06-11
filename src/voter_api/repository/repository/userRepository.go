package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"voter_api/connections"
	domain "voter_api/domain/user"
)

func RegisterUser(user *domain.User) error {
	//client := connections.GetInstanceMongoClient()
	//usersDatabase := client.Database("users")
	//uruguayVotersCollection := usersDatabase.Collection("uruguayVoters")
	//_, err2 := uruguayVotersCollection.InsertOne(context.TODO(), bson.M{"id": user.Id, "username": user.Username, "password": user.HashedPassword, "role": user.Role})
	//if err2 != nil {
	//	fmt.Println("error creating user")
	//	if err2 == mongo.ErrNoDocuments {
	//		return nil
	//	}
	//	log.Fatal(err2)
	//}
	//return err2
	return nil
}

func FindVoter(idVoter string) (*domain.User, error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_election")
	uruguayCollection := votesDatabase.Collection("voters")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", idVoter}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no voter habilitated to vote with that id")
		if err2 == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err2)
	}
	user := &domain.User{
		Id:        result["id"].(string),
		IdCircuit: result["idCircuit"].(string),
		Name:      result["name"].(string),
		LastName:  result["lastName"].(string),
		Voted:     int(result["voted"].(int32)),
		////Role:           result["role"].(string),
		//HashedPassword: result["password"].(string),
	}
	return user, nil
}

func FindCandidate(idCandidate string) (string, error) {
	client := connections.GetInstanceMongoClient()
	votesDatabase := client.Database("uruguay_votes")
	uruguayCollection := votesDatabase.Collection("votes")
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

	uruguayCollection = uruguayDataBase.Collection("total_votes")
	filter = bson.D{{"id", 1}}
	update = bson.D{{"$inc", bson.D{{"votes_counted", 1}}}}
	_, err2 = uruguayCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		fmt.Println("error registering new vote for candidate")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
