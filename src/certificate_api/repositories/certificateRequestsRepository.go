package repositories

import (
	"certificate_api/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	collection = "certificate_requests"
)

type requestsRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewRequestsRepo(mongoClient *mongo.Client, database string) *requestsRepo {
	return &requestsRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (reqRepo *requestsRepo) StoreRequest(req *models.CertificateRequestModel) error {
	certificatesDatabase := reqRepo.mongoClient.Database(reqRepo.database)
	requestsCollection := certificatesDatabase.Collection(collection)
	_, err2 := requestsCollection.InsertOne(context.TODO(), req)
	if err2 != nil {
		fmt.Println("error storing request")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
