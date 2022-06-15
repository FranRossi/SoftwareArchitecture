package repositories

import (
	m "consulting_api/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
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
	votesPerCandidates := result["votes_per_candidates"].(bson.A)
	voterPerParties := result["votes_per_parties"].(bson.A)
	regions := result["regions"].(bson.A)
	votesPerCandidatesStruct, votesPerPartiesStruct, regionsStruct := convertInterfaceResultToStruct(votesPerCandidates, voterPerParties, regions)

	resultElection := m.ResultElection{
		ElectionId:          electionId,
		AmountOfVotes:       int(result["amount_voted"].(int32)),
		TotalAmountOfVoters: int(result["total_amount_of_voters"].(int32)),
		VotesPerCandidates:  votesPerCandidatesStruct,
		VotesPerParties:     votesPerPartiesStruct,
		Regions:             regionsStruct,
	}
	return resultElection, nil
}

func convertInterfaceResultToStruct(votesPerCandidates, votesPerParties, regions bson.A) ([]m.CandidateEssential, []m.PoliticalPartyEssentials, []m.Region) {
	var votesPerCandidatesStruct []m.CandidateEssential
	var votesPerPartiesStruct []m.PoliticalPartyEssentials
	var regionsStruct []m.Region
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		votesPerCandidatesStruct = convertVotesPerCandidateToStruct(votesPerCandidates)
		wg.Done()
	}()
	go func() {
		votesPerPartiesStruct = convertVotesPerPartiesToStruct(votesPerParties)
		wg.Done()
	}()
	go func() {
		regionsStruct = convertRegionsToStruct(regions)
		wg.Done()
	}()
	wg.Wait()
	return votesPerCandidatesStruct, votesPerPartiesStruct, regionsStruct
}

func convertVotesPerCandidateToStruct(votesPerCandidates bson.A) []m.CandidateEssential {
	var votesPerCandidatesStruct []m.CandidateEssential
	for _, votePerCandidate := range votesPerCandidates {
		vote := votePerCandidate.(bson.M)
		votesPerCandidatesStruct = append(votesPerCandidatesStruct, m.CandidateEssential{
			Id:             vote["id"].(string),
			Name:           vote["name"].(string),
			Votes:          int(vote["votes"].(int32)),
			PoliticalParty: vote["politicalParty"].(string),
		})
	}
	return votesPerCandidatesStruct
}

func convertVotesPerPartiesToStruct(votesPerParties bson.A) []m.PoliticalPartyEssentials {
	var votesPerPartiesStruct []m.PoliticalPartyEssentials
	for _, votePerParty := range votesPerParties {
		vote := votePerParty.(bson.M)
		votesPerPartiesStruct = append(votesPerPartiesStruct, m.PoliticalPartyEssentials{
			Name:  vote["name"].(string),
			Votes: int(vote["votes"].(int32)),
		})
	}
	return votesPerPartiesStruct
}

func convertRegionsToStruct(regions bson.A) []m.Region {
	var regionsStruct []m.Region
	for _, region := range regions {
		r := region.(bson.M)
		regionsStruct = append(regionsStruct, m.Region{
			Name:        r["name"].(string),
			TotalVoters: int(r["total_voters"].(int32)),
			Votes:       int(r["votes"].(int32)),
		})
	}
	return regionsStruct
}
