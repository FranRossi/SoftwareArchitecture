package repository

import (
	"auth/connections"
	"auth/models"
	"context"
	"errors"
	"log"
	"os"
	l "own_logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	uruguayUsersCollection := usersDatabase.Collection(os.Getenv("COLLECTION_NAME"))
	_, err2 := uruguayUsersCollection.InsertOne(context.TODO(), user)
	if err2 != nil {
		l.LogError(err2.Error())
		if err2 == mongo.ErrNoDocuments {
			return errors.New("error creating user")
		}
		log.Fatal(err2)
	}
	return err2
}

func (repo *UsersRepo) FindUser(idUser string) (*models.UserDB, error) {
	client := connections.GetInstanceMongoClient()
	usersDatabase := client.Database(repo.database)
	uruguayUsersCollection := usersDatabase.Collection(os.Getenv("COLLECTION_NAME"))
	var result bson.M
	err2 := uruguayUsersCollection.FindOne(context.TODO(), bson.D{{"id", idUser}}).Decode(&result)
	if err2 != nil {
		l.LogWarning(err2.Error())
		if err2 == mongo.ErrNoDocuments {
			return nil, errors.New("there is no user with that id")
		}
	}
	user := &models.UserDB{
		Id:             result["id"].(string),
		Role:           result["role"].(string),
		HashedPassword: result["hashedpassword"].(string),
	}
	return user, nil
}
