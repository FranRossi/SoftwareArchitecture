package UruguayanElection

import (
	"ElectoralService/UruguayanElection/models"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func CreateElectionModel(amountVoters int) {
	CustomGeneratorVoter(amountVoters)
	CustomGeneratorPoliticalParties()
	CustomGeneratorCandidates()
	CustomGeneratorCircuits()
	electionModel := models.ElectionModel{}
	_ = faker.AddProvider("customIdFaker", func(v reflect.Value) (interface{}, error) {
		return strconv.FormatInt(int64(1), 10), nil
	})
	err := faker.FakeData(&electionModel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", electionModel)
}

func CustomGeneratorVoter(amountVoters int) {
	_ = faker.AddProvider("custom_voter", func(v reflect.Value) (interface{}, error) {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		departments := []string{"Montevideo", "Colonia", "Rocha", "Florida", "San Jose", "Soriano", "Salto", "Paysandu", "Treinta y Tres", "Artigas"}
		voters := make([]models.VoterModel, 0, 100000)
		for i := 0; i < amountVoters; i++ {
			voterModel := models.VoterModel{}
			voterModel.Id = strconv.FormatInt(int64(i+10000000), 10)
			departamentNumber := rand.Intn(len(departments))
			voterModel.Department = departments[departamentNumber]
			// TODO
			//voterModel.CivicCredential
			voterModel.IdCircuit = strconv.FormatInt(int64(departamentNumber), 10)
			faker.FakeData(&voterModel)
			voters = append(voters, voterModel)

		}
		return voters, nil
	})
}

func CustomGeneratorPoliticalParties() {
	_ = faker.AddProvider("custom_political_party", func(v reflect.Value) (interface{}, error) {
		politicalParties := make([]models.PoliticalPartyModel, 0, 4)
		partidoNacional := models.PoliticalPartyModel{
			Id:   "1",
			Name: "Partido Nacional",
		}
		partidoColorado := models.PoliticalPartyModel{
			Id:   "2",
			Name: "Partido Colorado",
		}
		politicalParties = append(politicalParties, partidoNacional, partidoColorado)
		return politicalParties, nil
	})
}

func CustomGeneratorCandidates() {
	_ = faker.AddProvider("custom_candidates", func(v reflect.Value) (interface{}, error) {
		candidates := make([]models.CandidateModel, 0, 4)
		candidate1NationalParty := models.CandidateModel{
			Id:               "1",
			IdPoliticalParty: "1",
		}
		candidate2NationalParty := models.CandidateModel{
			Id:               "2",
			IdPoliticalParty: "1",
		}
		candidateRedParty := models.CandidateModel{
			Id:               "2",
			IdPoliticalParty: "2",
		}
		faker.FakeData(&candidate1NationalParty)
		faker.FakeData(&candidate2NationalParty)
		faker.FakeData(&candidateRedParty)
		candidates = append(candidates, candidate1NationalParty, candidate2NationalParty, candidateRedParty)
		return candidates, nil
	})
}

func CustomGeneratorCircuits() {
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
