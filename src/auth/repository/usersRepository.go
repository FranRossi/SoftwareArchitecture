package repository

import (
	"auth/connections"
	"auth/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	Collection = "users"
)

type UsersRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewUsersRepo(mongoClient *mongo.Client, database string) *UsersRepo {
	return &UsersRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (repo *UsersRepo) RegisterUser(user *models.UserDB) error {
	client := connections.GetInstanceMongoClient()
	usersDatabase := client.Database(repo.database)
	uruguayUsersCollection := usersDatabase.Collection(Collection)
	_, err2 := uruguayUsersCollection.InsertOne(context.TODO(), user)
	if err2 != nil {
		fmt.Println("error creating user")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return err2
}

func (repo *UsersRepo) FindUser(idUser string) (*models.UserDB, error) {
	client := connections.GetInstanceMongoClient()
	usersDatabase := client.Database(repo.database)
	uruguayUsersCollection := usersDatabase.Collection(Collection)
	var result bson.M
	err2 := uruguayUsersCollection.FindOne(context.TODO(), bson.D{{"id", idUser}}).Decode(&result)
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println("there is no user with that id")
		if err2 == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err2)
	}
	user := &models.UserDB{
		Id:             result["id"].(string),
		Role:           result["role"].(string),
		HashedPassword: result["hashedpassword"].(string),
	}
	return user, nil
}
