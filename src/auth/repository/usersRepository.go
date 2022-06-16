package repositories

import (
	"auth/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Collection = "users"
)

type UsersRepo struct {
	mongoClient *mongo.Client
	database    string
}


