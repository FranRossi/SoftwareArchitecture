package repository

import (
	"context"
	"electoral_service/connections"
	"electoral_service/models"
	"encoding/json"
	"encrypt"
	"fmt"
	mq "message_queue"
	"os"
	l "own_logger"

	"go.mongodb.org/mongo-driver/bson"
)

type ElectionRepo struct {
}

func DropDataBases() {
	client := connections.GetInstanceMongoClient()
	electionDataBase := client.Database(os.Getenv("ELECTION_DB"))
	electionDataBase.Drop(context.TODO())
	electionDataBase = client.Database(os.Getenv("VOTES_DB"))
	electionDataBase.Drop(context.TODO())
}

func (repo *ElectionRepo) StoreElectionConfiguration(election *models.ElectionModelEssential) error {
	client := connections.GetInstanceMongoClient()
	electionDataBase := client.Database(os.Getenv("ELECTION_DB"))
	electionCollection := electionDataBase.Collection(os.Getenv("CONF_COL"))
	_, err := electionCollection.InsertOne(context.TODO(), bson.M{"id": election.Id, "startingDate": election.StartingDate, "finishingDate": election.FinishingDate, "electionMode": election.ElectionMode, "otherField": election.OtherFields})
	if err != nil {
		return fmt.Errorf("error storing election configuration")
	}

	votesDataBase := client.Database(os.Getenv("VOTES_DB"))
	collectionTotalVotes := votesDataBase.Collection(os.Getenv("TOTAL_VOTES_COL"))
	amountVotes := bson.D{{"votes_counted", 0}, {"election_id", election.Id}}
	_, err = collectionTotalVotes.InsertOne(context.TODO(), amountVotes)
	if err != nil {
		return fmt.Errorf("error storing initial amount of votes")
	}
	return err
}

func StoreElectionVoters(electionId string, voters []models.VoterModel) error {

	// Don't use for voter := range voters, because it won't change the properties of the voters in the original array
	for i := range voters {
		sendStatsOfVoterToStatsService(electionId, &voters[i])
		encrypt.EncryptVoter(&voters[i])
	}
	votersInterface := convertVotersModelToInterface(voters)
	client := connections.GetInstanceMongoClient()
	electionDataBase := client.Database(os.Getenv("ELECTION_DB"))
	votersCollection := electionDataBase.Collection(os.Getenv("VOTERS_COL"))
	_, err := votersCollection.InsertMany(context.TODO(), votersInterface)
	if err != nil {
		return fmt.Errorf("error storing voters")
	}
	return nil
}

func sendStatsOfVoterToStatsService(electionId string, voter *models.VoterModel) {
	type VoterStats struct {
		ElectionId string
		BirthDate  string
		Region     string
		Circuit    string
		Sex        string
	}

	var voterStats VoterStats
	voterStats.BirthDate = voter.BirthDate
	voterStats.Circuit = voter.OtherFields["circuit"].(string)
	voterStats.Region = voter.Region
	voterStats.Sex = voter.Sex
	voterStats.ElectionId = electionId

	jsonStats, errs := json.Marshal(voterStats)
	if errs != nil {
		l.LogError("error sending voter stats to queue:" + errs.Error())
	}
	mq.GetMQWorker().Send("stats-total", jsonStats)
}

func convertVotersModelToInterface(voters []models.VoterModel) []interface{} {
	var votersInterface []interface{}
	for _, v := range voters {
		votersInterface = append(votersInterface, v)
	}
	return votersInterface
}

func StoreCandidates(candidates []models.CandidateModel) error {
	client := connections.GetInstanceMongoClient()
	candidatesToStore := convertCandidateModelToInterface(candidates)
	votesDataBase := client.Database(os.Getenv("VOTES_DB"))
	candidatesCollection := votesDataBase.Collection(os.Getenv("CANDIDATES_VOTES_COL"))
	_, err := candidatesCollection.InsertMany(context.TODO(), candidatesToStore)
	if err != nil {
		return fmt.Errorf("error storing initial candidates")
	}
	return nil
}

func convertCandidateModelToInterface(candidates []models.CandidateModel) []interface{} {
	var candidatesResume []models.CandidateEssential
	for _, v := range candidates {
		candidatesResume = append(candidatesResume, models.CandidateEssential{Id: v.Id, Name: v.FullName, Votes: 0, PoliticalParty: v.NamePoliticalParty})
	}
	var candidatesInterface []interface{}
	for _, v := range candidatesResume {
		candidatesInterface = append(candidatesInterface, v)
	}
	return candidatesInterface
}

func GetTotalVotes(electionId string) (int, error) {
	client := connections.GetInstanceMongoClient()
	votesDataBase := client.Database(os.Getenv("VOTES_DB"))
	votesCollection := votesDataBase.Collection(os.Getenv("TOTAL_VOTES_COL"))
	var amountVotes bson.M
	err := votesCollection.FindOne(context.TODO(), bson.D{{"election_id", electionId}}).Decode(&amountVotes)
	if err != nil {
		return -1, fmt.Errorf("error getting amount of votes")
	}
	amountVotesCounted := int(amountVotes["votes_counted"].(int32))
	return amountVotesCounted, nil
}

func GetEachCandidatesVotes() ([]models.CandidateEssential, error) {
	client := connections.GetInstanceMongoClient()
	votesDataBase := client.Database(os.Getenv("VOTES_DB"))
	votesCollection := votesDataBase.Collection(os.Getenv("CANDIDATES_VOTES_COL"))
	var results []bson.M
	cursor, err := votesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return []models.CandidateEssential{}, fmt.Errorf("there are no votes: %v", err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []models.CandidateEssential{}, err
	}
	var votesCandidates []models.CandidateEssential
	for _, result := range results {
		candidate := &models.CandidateEssential{
			Votes:          int(result["votes"].(int32)),
			Name:           result["name"].(string),
			Id:             result["id"].(string),
			PoliticalParty: result["politicalParty"].(string),
		}
		votesCandidates = append(votesCandidates, *candidate)
	}
	return votesCandidates, nil
}

func StoreElectionResult(resultElection models.ResultElection) error {
	client := connections.GetInstanceMongoClient()
	votesDataBase := client.Database(os.Getenv("VOTES_DB"))
	votesCollection := votesDataBase.Collection(os.Getenv("RESULT_COL"))
	_, err := votesCollection.InsertOne(context.TODO(), resultElection)
	if err != nil {
		message := "error storing result election"
		return fmt.Errorf(message)
	}
	return nil
}
