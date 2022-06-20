package repositories

import (
	"cache"
	m "consulting_api/models"
	"context"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VotesRepo struct {
	mongoClient *mongo.Client
	database    string
}

func NewRequestsRepo(mongoClient *mongo.Client, database string) *VotesRepo {
	return &VotesRepo{
		mongoClient: mongoClient,
		database:    database,
	}
}

func (certRepo *VotesRepo) RequestVote(voterId, electionId string) (*m.VoteModel, error) {

	var voteFromCache m.VoteModel
	errCache := cache.Get(voteCacheKey(voterId, electionId), &voteFromCache)
	if errCache == nil {
		return &voteFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("votes_info")
	filter := bson.D{{"voter", voterId}, {"election", electionId}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return &m.VoteModel{}, fmt.Errorf("error requesting voter on election: %v", err)
	}
	vote := &m.VoteModel{
		VoterId:   result["voter"].(string),
		TimeVoted: result["time_front_end"].(string),
	}
	cache.Save(voteCacheKey(voterId, electionId), vote, cache.DefaultExpiration)
	return vote, nil
}
func voteCacheKey(voterId, electionId string) string {
	return "election_" + electionId + "_voter_" + voterId
}

func (certRepo *VotesRepo) RequestElectionResult(electionId string) (m.ResultElection, error) {

	var resultFromCache m.ResultElection
	errCache := cache.Get(electionResultCachePrefix+electionId, &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("result_election")
	filter := bson.D{{"election_id", electionId}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return m.ResultElection{}, fmt.Errorf("error requesting election result: %v", err)
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
	cache.Save(electionResultCachePrefix+electionId, resultElection, cache.DefaultExpiration)
	return resultElection, nil
}

const electionResultCachePrefix = "election_result_"

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

func (certRepo *VotesRepo) RequestPopularVotingTimes(electionId string) (map[string]int, error) {

	var resultFromCache map[string]int
	errCache := cache.Get(popularTimeCachePrefix+electionId, &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("votes_info")
	filter := bson.D{{"election", electionId}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return map[string]int{}, fmt.Errorf("error requesting popular voting times: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return map[string]int{}, fmt.Errorf("error requesting popular voting times using cursor: %v", err)
	}
	averageTimes := calculateWhichHoursHaveMoreVotes(results)
	cache.Save(popularTimeCachePrefix+electionId, averageTimes, cache.DefaultExpiration)
	return averageTimes, nil
}

const popularTimeCachePrefix = "popular_time__election_"

func calculateWhichHoursHaveMoreVotes(results []bson.M) map[string]int {
	hours := make(map[string]int)
	for _, result := range results {
		timeVoted := result["time_front_end"].(string)
		timeVotedSplit := strings.Split(timeVoted, "T")
		hourMinutesSecond := timeVotedSplit[1]
		hourSplit := strings.Split(hourMinutesSecond, ":")
		hour := hourSplit[0]
		if _, ok := hours[hour]; ok {
			hours[hour]++
		} else {
			hours[hour] = 1
		}
	}
	return hours
}

const votesPerCircuitCachePrefix = "votes_per_circuit_election_"

func (certRepo *VotesRepo) RequestVotesPerCircuits(electionId string) ([]m.VotesPerCircuits, error) {
	var resultFromCache []m.VotesPerCircuits
	errCache := cache.Get(votesPerCircuitCachePrefix+electionId, &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("statistics")
	filter := bson.D{{"election_id", electionId}, {"group_type", "circuit"}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return []m.VotesPerCircuits{}, fmt.Errorf("error requesting votes per circuits: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []m.VotesPerCircuits{}, fmt.Errorf("error requesting votes per circuits using cursor: %v", err)
	}
	votesPerCircuits := convertVotesPerCircuitsToStruct(electionId, results)
	cache.Save(votesPerCircuitCachePrefix+electionId, votesPerCircuits, cache.DefaultExpiration)
	return votesPerCircuits, nil
}

func convertVotesPerCircuitsToStruct(electionId string, results []bson.M) []m.VotesPerCircuits {
	var votesPerCircuits []m.VotesPerCircuits
	for _, result := range results {
		perCircuit := m.VotesPerCircuits{
			ElectionId:       electionId,
			Circuit:          result["circuit"].(string),
			VotesPerCircuits: int(result["votes"].(int32)),
			GroupName:        result["group_name"].(string),
		}
		votesPerCircuits = append(votesPerCircuits, perCircuit)
	}
	return votesPerCircuits
}

const votesPerRegionCachePrefix = "votes_per_region_election_"

func (certRepo *VotesRepo) RequestVotesPerRegions(electionId string) ([]m.VotesPerRegion, error) {
	var resultFromCache []m.VotesPerRegion
	errCache := cache.Get(votesPerRegionCachePrefix+electionId, &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("statistics")
	filter := bson.D{{"election_id", electionId}, {"group_type", "region"}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return []m.VotesPerRegion{}, fmt.Errorf("error requesting votes per regions: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []m.VotesPerRegion{}, fmt.Errorf("error requesting votes per regions using cursor: %v", err)
	}
	votesPerRegions := convertVotesPerRegionsToStruct(electionId, results)
	cache.Save(votesPerRegionCachePrefix+electionId, votesPerRegions, cache.DefaultExpiration)
	return votesPerRegions, nil
}

func convertVotesPerRegionsToStruct(electionId string, results []bson.M) []m.VotesPerRegion {
	var votesPerCircuits []m.VotesPerRegion
	for _, result := range results {
		perCircuit := m.VotesPerRegion{
			ElectionId:     electionId,
			Region:         result["region"].(string),
			VotesPerRegion: int(result["votes"].(int32)),
			GroupName:      result["group_name"].(string),
		}
		votesPerCircuits = append(votesPerCircuits, perCircuit)
	}
	return votesPerCircuits
}
