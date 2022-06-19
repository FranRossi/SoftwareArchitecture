package repositories

import (
	m "consulting_api/models"
	"context"
	"fmt"
	l "own_logger"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ElectionRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewElectionRepo(mongoClient *mongo.Client, database string) *ElectionRepo {
	return &ElectionRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (certRepo *ElectionRepo) RequestElectionConfig(electionId string) (m.ElectionConfig, error) {

	var configs m.ElectionConfig
	configFromCache, errCache := RequestElectionConfigFromCache(electionId)
	if errCache == nil {
		l.LogInfo("Election config for election " + electionId + " was found in cache")
		return configFromCache, nil
	} else {
		defer SaveElectionConfigToCache(configs)
	}

	client := certRepo.mongoClient
	electionDatabase := client.Database("uruguay_election")
	uruguayCollection := electionDatabase.Collection("configuration_election")
	var result bson.M
	err2 := uruguayCollection.FindOne(context.TODO(), bson.D{{"id", electionId}}).Decode(&result)
	if err2 != nil {
		return m.ElectionConfig{}, fmt.Errorf("election not found: %v", err2)
	}
	maxVotesString := result["otherField"].(bson.M)["maxVotes"].(string)
	maxVotes, err3 := strconv.Atoi(maxVotesString)
	maxCertificatesString := result["otherField"].(bson.M)["maxCertificate"].(string)
	maxCertificates, err3 := strconv.Atoi(maxCertificatesString)
	emails := result["otherField"].(bson.M)["emails"].(bson.A)
	var emailsArray []string
	for _, email := range emails {
		emailsArray = append(emailsArray, email.(string))
	}
	if err3 != nil {
		return m.ElectionConfig{}, fmt.Errorf("worng maximum values: %v", err3)
	}
	configs = m.ElectionConfig{
		MaxVotes:        maxVotes,
		MaxCertificates: maxCertificates,
		Emails:          emailsArray,
	}
	return configs, nil
}
