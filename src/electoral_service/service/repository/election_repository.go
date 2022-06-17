package repository

import (
	"context"
	"electoral_service/connections"
	"electoral_service/models"
	"encrypt"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ElectionRepo struct {
}

func DropDataBases() {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayDataBase.Drop(context.TODO())
	uruguayDataBase = client.Database("uruguay_votes")
	uruguayDataBase.Drop(context.TODO())
}

func (repo *ElectionRepo) StoreElectionConfiguration(election *models.ElectionModelEssential) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayanElectionCollection := uruguayDataBase.Collection("configuration_election")
	_, err := uruguayanElectionCollection.InsertOne(context.TODO(), bson.M{"id": election.Id, "startingDate": election.StartingDate, "finishingDate": election.FinishingDate, "electionMode": election.ElectionMode, "otherField": election.OtherFields})
	if err != nil {
		fmt.Println("error storing election configuration")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}

	uruguayVotesDataBase := client.Database("uruguay_votes")
	collectionTotalVotes := uruguayVotesDataBase.Collection("total_votes")
	amountVotes := bson.D{{"votes_counted", 0}, {"election_id", election.Id}}
	_, err = collectionTotalVotes.InsertOne(context.TODO(), amountVotes)
	if err != nil {
		fmt.Println("error storing initial amount of votes")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return err
}

func StoreElectionVoters(voters []models.VoterModel) error {

	encryptVoter(voters)
	votersInterface := convertVotersModelToInterface(voters)
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_election")
	uruguayanVotersCollection := uruguayDataBase.Collection("voters")
	_, err := uruguayanVotersCollection.InsertMany(context.TODO(), votersInterface)
	if err != nil {
		fmt.Println("error storing voters")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return err
}

func encryptVoter(voters []models.VoterModel) {

	// Don't use for voter := range voters, because it won't change the properties of the voters in the original array
	for i := range voters {
		encrypt.EncryptVoter(&voters[i])
	}
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
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanCandidatesCollection := uruguayDataBase.Collection("votes_per_candidate")
	_, err := uruguayanCandidatesCollection.InsertMany(context.TODO(), candidatesToStore)
	if err != nil {
		fmt.Println("error storing initial candidates")
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
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
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanVotesCollection := uruguayDataBase.Collection("total_votes")
	var amountVotes bson.M
	err := uruguayanVotesCollection.FindOne(context.TODO(), bson.D{{"election_id", electionId}}).Decode(&amountVotes)
	if err != nil {
		fmt.Println("error getting amount of votes")
		if err == mongo.ErrNoDocuments {
			return -1, nil
		}
		log.Fatal(err)
	}
	amountVotesCounted := int(amountVotes["votes_counted"].(int32))
	return amountVotesCounted, nil
}

func GetEachCandidatesVotes() ([]models.CandidateEssential, error) {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanVotesCollection := uruguayDataBase.Collection("votes_per_candidate")
	var results []bson.M
	cursor, err := uruguayanVotesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("there are no votes")
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
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
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanVotesCollection := uruguayDataBase.Collection("result_election")
	_, err := uruguayanVotesCollection.InsertOne(context.TODO(), resultElection)
	if err != nil {
		message := "error storing result election"
		return fmt.Errorf(message)
	}
	return nil
}
