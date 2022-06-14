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

func (certRepo *ConsultingRepo) RequestElectionResult(electionId string) (m.ResultElection, error) {
	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("result_election")
	filter := bson.D{{"election_id", electionId}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return m.ResultElection{}, err
	}
	resultElection := m.ResultElection{
		ElectionId:          electionId,
		AmountOfVotes:       int(result["amount_voted"].(int32)),
		TotalAmountOfVoters: int(result["total_amount_voters"].(int32)),
		VotesPerCandidates:  result["votes_per_candidate"].([]m.CandidateEssential),
		VotesPerParties:     result["votes_per_party"].([]m.PoliticalPartyEssentials),
	}
	return resultElection, nil
}

func (certRepo *ConsultingRepo) RequestElectionResultPerDepartment(electionId string) (int, int) {
	//client := certRepo.mongoClient
	//database := client.Database(certRepo.database)
	return 1, 1
}
