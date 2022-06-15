package controller

import (
	"electoral_service/adapter/uruguayan_election"
	models2 "electoral_service/adapter/uruguayan_election/models"
	modelsGeneric "electoral_service/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ElectionController struct {
}

func (controller *ElectionController) GetElectionSettings() *modelsGeneric.ElectionModelEssential {
	uruguayan_election.ConfigEnvironment()
	url := os.Getenv("electoral_service_url")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var electionSettings models2.ElectionJson
	er := json.Unmarshal(jsonBytes, &electionSettings)
	if er != nil {
		log.Fatal(er)
	}
	if electionSettings.Error {
		log.Fatal(electionSettings.Msg)
	}
	return convertToElectionModel(electionSettings.Election)

}

func convertToElectionModel(election models2.ElectionModel) *modelsGeneric.ElectionModelEssential {
	electionGeneric := &modelsGeneric.ElectionModelEssential{
		Id:               election.Id,
		StartingDate:     election.StartingDate,
		FinishingDate:    election.FinishingDate,
		ElectionMode:     election.ElectionMode,
		Voters:           convertVoters(election.Voters),
		PoliticalParties: convertPoliticalParties(election.PoliticalParties),
		OtherFields:      convertMaximumVotesAndCertificate(),
	}
	return electionGeneric
}

func convertMaximumVotesAndCertificate() map[string]any {
	maximum := map[string]any{
		"maxVotes":       os.Getenv("maxVotes"),
		"maxCertificate": os.Getenv("maxCertificate"),
	}
	return maximum

}

func convertVoters(voters []models2.VoterModel) []modelsGeneric.VoterModel {
	var votersGeneric []modelsGeneric.VoterModel
	for _, voter := range voters {
		voterGeneric := modelsGeneric.VoterModel{
			Id:          voter.Id,
			FullName:    voter.Name + " " + voter.LastName,
			BirthDate:   voter.BirthDate,
			Email:       voter.Email,
			Sex:         voter.Sex,
			Phone:       voter.Phone,
			Voted:       voter.Voted,
			Region:      voter.Department,
			OtherFields: map[string]any{"lastname": voter.LastName, "civiccredential": voter.CivicCredential, "circuit": voter.IdCircuit},
		}
		votersGeneric = append(votersGeneric, voterGeneric)
	}
	return votersGeneric
}

func convertPoliticalParties(politicalParties []models2.PoliticalPartyModel) []modelsGeneric.PoliticalPartyModel {
	var politicalPartiesGeneric []modelsGeneric.PoliticalPartyModel
	for _, politicalParty := range politicalParties {
		politicalPartyGeneric := modelsGeneric.PoliticalPartyModel{
			Id:          politicalParty.Id,
			Name:        politicalParty.Name,
			Candidates:  convertCandidates(politicalParty.Candidates, politicalParty.Name),
			OtherFields: map[string]any{},
		}
		politicalPartiesGeneric = append(politicalPartiesGeneric, politicalPartyGeneric)
	}
	return politicalPartiesGeneric
}

func convertCandidates(candidates []models2.CandidateModel, politicalParty string) []modelsGeneric.CandidateModel {
	var candidatesGeneric []modelsGeneric.CandidateModel
	for _, candidate := range candidates {
		candidateGeneric := modelsGeneric.CandidateModel{
			Id:                 candidate.Id,
			FullName:           candidate.Name + " " + candidate.LastName,
			Sex:                candidate.Sex,
			BirthDate:          candidate.BirthDate,
			IdPoliticalParty:   candidate.IdPoliticalParty,
			NamePoliticalParty: politicalParty,
			OtherFields:        map[string]any{"lastname": candidate.LastName},
		}
		candidatesGeneric = append(candidatesGeneric, candidateGeneric)
	}
	return candidatesGeneric
}
