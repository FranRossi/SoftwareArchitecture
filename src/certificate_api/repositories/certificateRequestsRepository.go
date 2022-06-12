package repositories

import (
	"certificate_api/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"certificate_api/connections"
)

func StoreRequest(req *models.CertificateRequestModel) error {
	client := connections.GetInstanceMongoClient()
	certificatesDatabase := client.Database("certificates")
	requestsCollection := certificatesDatabase.Collection("certificate_requests")
	_, err2 := requestsCollection.InsertOne(context.TODO(), req)
	if err2 != nil {
		fmt.Println("error storing vote")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
