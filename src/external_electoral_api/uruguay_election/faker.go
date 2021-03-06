package uruguay_election

import (
	"external_electoral_api/uruguay_election/models"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func CreateElectionMock(id string, amountVoters int) (models.ElectionModel, error) {
	customGeneratorVoter(amountVoters)
	candidates := customGeneratorCandidates()
	customGeneratorPoliticalParties(candidates)
	customGeneratorCircuits()
	electionModel := models.ElectionModel{}
	_ = faker.AddProvider("customIdFaker", func(v reflect.Value) (interface{}, error) {
		return id, nil
	})
	customGenerateDates(&electionModel)
	err := faker.FakeData(&electionModel)
	if err != nil {
		return models.ElectionModel{}, err
	}
	setEmail(&electionModel)
	return electionModel, err
}
func setEmail(election *models.ElectionModel) {
	election.Emails = []string{"montevideo@intendencia.com", "montevideo@presidencia.com", "colonia@intendencia.com"}
}

func customGenerateDates(election *models.ElectionModel) {
	election.StartingDate = time.Now().Format(time.RFC3339)
	election.FinishingDate = time.Now().Add(time.Minute * 59).Format(time.RFC3339)
	//election.FinishingDate = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
}

func customGeneratorVoter(amountVoters int) {
	_ = faker.AddProvider("custom_voter", func(v reflect.Value) (interface{}, error) {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		departments := []string{"Montevideo", "Colonia", "Rocha", "Florida", "San Jose", "Soriano", "Salto", "Paysandu", "Treinta y Tres", "Artigas"}
		voters := make([]models.VoterModel, 0, amountVoters)
		for i := 0; i < amountVoters; i++ {
			voterModel := models.VoterModel{}
			voterModel.Id = strconv.FormatInt(int64(i+10000000), 10)
			departmentNumber := rand.Intn(len(departments))
			voterModel.Department = departments[departmentNumber]
			if voterModel.Id == "10000000" {
				voterModel.IdCircuit = "1"
			} else {
				voterModel.IdCircuit = strconv.FormatInt(int64(departmentNumber), 10)
			}
			voterModel.CivicCredential = randomCivicCredential()
			faker.FakeData(&voterModel)
			voters = append(voters, voterModel)

		}
		return voters, nil
	})
}

func randomCivicCredential() string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numbersBytes = "0123456789"
	b := make([]byte, 8)
	for i := range b {
		if i < 3 {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else {
			b[i] = numbersBytes[rand.Intn(len(numbersBytes))]
		}
	}
	return string(b)
}

func customGeneratorPoliticalParties(candidates [][]models.CandidateModel) {
	_ = faker.AddProvider("custom_political_party", func(v reflect.Value) (interface{}, error) {
		politicalParties := make([]models.PoliticalPartyModel, 0, 4)
		partidoNacional := models.PoliticalPartyModel{
			Id:         "1",
			Name:       "Partido Gaviota",
			Candidates: candidates[0],
		}
		partidoColorado := models.PoliticalPartyModel{
			Id:         "2",
			Name:       "Partido Orangutan",
			Candidates: candidates[1],
		}
		politicalParties = append(politicalParties, partidoNacional, partidoColorado)
		return politicalParties, nil
	})
}

func customGeneratorCandidates() [][]models.CandidateModel {
	candidates := make([][]models.CandidateModel, 0, 4)
	candidatesRed := make([]models.CandidateModel, 0, 2)
	candidatesBlue := make([]models.CandidateModel, 0, 3)
	candidate1NationalParty := models.CandidateModel{
		Id:               "1",
		IdPoliticalParty: "1",
	}
	candidate2NationalParty := models.CandidateModel{
		Id:               "2",
		IdPoliticalParty: "1",
	}
	candidateRedParty := models.CandidateModel{
		Id:               "3",
		IdPoliticalParty: "2",
	}
	faker.FakeData(&candidate1NationalParty)
	faker.FakeData(&candidate2NationalParty)
	faker.FakeData(&candidateRedParty)
	candidatesRed = append(candidatesRed, candidateRedParty)
	candidatesBlue = append(candidatesBlue, candidate1NationalParty, candidate2NationalParty)
	candidates = append(candidates, candidatesBlue, candidatesRed)
	return candidates
}

func customGeneratorCircuits() {
	_ = faker.AddProvider("custom_circuits", func(v reflect.Value) (interface{}, error) {
		departments := []string{"Montevideo", "Colonia", "Rocha", "Florida", "San Jose", "Soriano", "Salto", "Paysandu", "Treinta y Tres", "Artigas"}
		circuits := make([]models.CircuitModel, 0, 10)
		for i := 0; i < 10; i++ {
			circuitModel := models.CircuitModel{}
			circuitModel.Id = strconv.FormatInt(int64(i), 10)
			circuitModel.Department = departments[i]
			faker.FakeData(&circuitModel)
			circuits = append(circuits, circuitModel)
		}
		return circuits, nil
	})
}
