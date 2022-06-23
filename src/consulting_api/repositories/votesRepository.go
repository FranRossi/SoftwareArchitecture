package repositories

import (
	"cache"
	m "consulting_api/models"
	"context"
	"fmt"
	"os"
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

func (certRepo *VotesRepo) RequestVote(voterId, electionId string) (m.VoteModel, error) {

	var voteFromCache m.VoteModel
	errCache := cache.Get(voteCacheKey(voterId, electionId), &voteFromCache)
	if errCache == nil {
		return voteFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(certRepo.database)
	collection := database.Collection("votes_info")
	filter := bson.D{{"voter", voterId}, {"election", electionId}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return m.VoteModel{}, fmt.Errorf("error requesting voter on election: %v", err)
	}
	vote := m.VoteModel{
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

func votesPerCircuitCacheKey(electionId string, circuit string) string {
	return "votes_per_circuit_election_" + electionId + "_" + circuit
}
func (certRepo *VotesRepo) RequestVotesPerCircuits(electionId string, circuit string) (m.VotesPerCircuits, error) {
	var resultFromCache m.VotesPerCircuits
	errCache := cache.Get(votesPerCircuitCacheKey(electionId, circuit), &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database("statistics")
	collection := database.Collection("votes_stats")
	filter := bson.D{{"election_id", electionId}, {"group_type", "circuit"}, {"circuit", circuit}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return m.VotesPerCircuits{}, fmt.Errorf("error requesting votes per circuits: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return m.VotesPerCircuits{}, fmt.Errorf("error requesting votes per circuits using cursor: %v", err)
	}
	votesPerGroup := convertVotesPerGroupToStruct(results)
	votesPerCircuits := m.VotesPerCircuits{
		ElectionId:   electionId,
		Circuit:      circuit,
		DataPerGroup: votesPerGroup,
	}
	cache.Save(votesPerCircuitCacheKey(electionId, circuit), votesPerCircuits, cache.DefaultExpiration)
	return votesPerCircuits, nil
}

func convertVotesPerGroupToStruct(results []bson.M) []m.VotesPerGroup {
	var votesPerCircuits []m.VotesPerGroup
	for _, result := range results {
		perCircuit := m.VotesPerGroup{

			GroupName:    result["group_name"].(string),
			MinAge:       int(result["min_age"].(int32)),
			MaxAge:       int(result["max_age"].(int32)),
			Sex:          result["sex"].(string),
			CurrentVotes: int(result["votes"].(int32)),
			Total:        int(result["capacity"].(int32)),
		}
		votesPerCircuits = append(votesPerCircuits, perCircuit)
	}
	return votesPerCircuits
}

func votesPerRegionCacheKey(electionId string, region string) string {
	return "votes_per_region_election_" + electionId + "_" + region
}

func (certRepo *VotesRepo) RequestVotesPerRegions(electionId string, region string) (m.VotesPerRegion, error) {
	var resultFromCache m.VotesPerRegion
	errCache := cache.Get(votesPerRegionCacheKey(electionId, region), &resultFromCache)
	if errCache == nil {
		return resultFromCache, nil
	}

	client := certRepo.mongoClient
	database := client.Database(os.Getenv("STATS_DB"))
	collection := database.Collection(os.Getenv("STATS_COLLECTION"))
	filter := bson.D{{"election_id", electionId}, {"group_type", "region"}, {"region", region}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return m.VotesPerRegion{}, fmt.Errorf("error requesting votes per regions: %v", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return m.VotesPerRegion{}, fmt.Errorf("error requesting votes per regions using cursor: %v", err)
	}
	votesPerGroup := convertVotesPerGroupToStruct(results)
	votesPerRegions := m.VotesPerRegion{
		ElectionId:   electionId,
		Region:       region,
		DataPerGroup: votesPerGroup,
	}
	cache.Save(votesPerRegionCacheKey(electionId, region), votesPerRegions, cache.DefaultExpiration)
	return votesPerRegions, nil
}
