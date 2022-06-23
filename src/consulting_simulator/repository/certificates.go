package repository

import (
	"consulting_simulator/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	certCollection = "certificates"
)

type CertificatesRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewRequestsRepo(mongoClient *mongo.Client, database string) *CertificatesRepo {
	return &CertificatesRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (certRepo *CertificatesRepo) FindAllCertificateForElection(electionId string) (req []models.CertificateResponseModel, err error) {
	certificatesDatabase := certRepo.mongoClient.Database(certRepo.database)
	requestsCollection := certificatesDatabase.Collection(certCollection)
	cursor, err2 := requestsCollection.Find(context.Background(), bson.D{{"id_election", electionId}})
	if err2 != nil {
		return []models.CertificateResponseModel{}, fmt.Errorf("there are no certificates: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []models.CertificateResponseModel{}, err
	}
	var certificates []models.CertificateResponseModel
	for _, result := range results {
		certificate := &models.CertificateResponseModel{
			IdVoter:            result["id_voter"].(string),
			IdElection:         result["id_election"].(string),
			TimeVoted:          result["time_voted"].(string),
			VoteIdentification: result["vote_identification"].(string),
			StartingDate:       result["starting_date"].(string),
			FinishingDate:      result["finishing_date"].(string),
			ElectionMode:       result["election_mode"].(string),
			Fullname:           result["fullname"].(string),
			Message:            result["message"].(string),
		}
		certificates = append(certificates, *certificate)
	}
	return certificates, nil
}
