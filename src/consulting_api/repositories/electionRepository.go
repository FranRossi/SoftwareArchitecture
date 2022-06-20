package repositories

import (
	"cache"
	m "consulting_api/models"
	"context"
	"fmt"
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

const electionConfigCacheKeyPrefix = "election_config_"

func (certRepo *ElectionRepo) RequestElectionConfig(electionId string) (m.ElectionConfig, error) {

	var configFromCache m.ElectionConfig
	errCache := cache.Get(electionConfigCacheKeyPrefix+electionId, &configFromCache)
	if errCache == nil {
		return configFromCache, nil
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
	configs := m.ElectionConfig{
		MaxVotes:        maxVotes,
		MaxCertificates: maxCertificates,
		Emails:          emailsArray,
	}
	cache.Save(electionConfigCacheKeyPrefix+electionId, &configs, cache.DefaultExpiration)
	return configs, nil
}
