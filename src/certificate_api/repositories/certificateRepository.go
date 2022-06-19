package repositories

import (
	"certificate_api/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	reqCollection  = "certificate_requests"
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

func (certRepo *CertificatesRepo) StoreRequest(req models.CertificateRequestModel) error {
	certificatesDatabase := certRepo.mongoClient.Database(certRepo.database)
	requestsCollection := certificatesDatabase.Collection(reqCollection)
	_, err2 := requestsCollection.InsertOne(context.TODO(), req)
	if err2 != nil {
		return fmt.Errorf("error storing request: %v", err2)
	}
	return nil
}

func (certRepo *CertificatesRepo) StoreCertificate(cert models.CertificateModel) error {
	certificatesDatabase := certRepo.mongoClient.Database(certRepo.database)
	certificatesCollection := certificatesDatabase.Collection(certCollection)
	_, err2 := certificatesCollection.InsertOne(context.TODO(), cert)
	if err2 != nil {
		return fmt.Errorf("error storing certificate: %v", err2)
	}
	return nil
}

func (certRepo *CertificatesRepo) FindCertificate(voterId, voteIdentification string) (req models.CertificateModel, err error) {
	certificatesDatabase := certRepo.mongoClient.Database(certRepo.database)
	requestsCollection := certificatesDatabase.Collection(certCollection)
	var result bson.M
	err2 := requestsCollection.FindOne(context.TODO(), bson.D{{"id", voterId}, {"voteIdentification", voteIdentification}}).Decode(&result)
	if err2 != nil {
		return req, fmt.Errorf("there is no certificate assigned to the voter with that id: %v", err2)
	}
	req = models.CertificateModel{
		IdVoter:            result["id_voter"].(string),
		IdElection:         result["id_election"].(string),
		TimeVoted:          result["time_voted"].(string),
		VoteIdentification: result["vote_identification"].(string),
		StartingDate:       result["starting_date"].(string),
		FinishingDate:      result["finishing_date"].(string),
		ElectionMode:       result["election_mode"].(string),
		Fullname:           result["fullname"].(string),
	}
	return req, nil
}
