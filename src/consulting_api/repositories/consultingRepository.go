package repositories

import (
	m "consulting_api/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConsultingRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewRequestsRepo(mongoClient *mongo.Client, database string) *ConsultingRepo {
	return &ConsultingRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (certRepo *ConsultingRepo) RequestVote(voterId, electionId string) (*m.VoteModel, error) {
	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("votes_info")
	filter := bson.D{{"voter", voterId}, {"election", electionId}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return &m.VoteModel{}, err
	}
	vote := &m.VoteModel{
		VoterId:   result["voter"].(string),
		TimeVoted: result["time_front_end"].(string),
	}
	return vote, nil
}
