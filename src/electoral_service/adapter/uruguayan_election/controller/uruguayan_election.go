package controller

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	modelsGeneric "electoral_service/models"
	"electoral_service/service/logic"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "http://localhost:8080/api/v1/election/uruguay/?id=1"

type ElectionController struct {
	electionLogic *logic.ElectionLogic
}

func NewElectionController(logic *logic.ElectionLogic) *ElectionController {
	return &ElectionController{electionLogic: logic}
}

func (controller *ElectionController) GetElectionSettings() *modelsGeneric.ElectionModelEssential {
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
	} else {
		return convertToElectionModel(electionSettings.Election)
	}
	return nil
}

func convertToElectionModel(election models2.ElectionModel) *modelsGeneric.ElectionModelEssential {
	electionGeneric := &modelsGeneric.ElectionModelEssential{
		Id:               election.Id,
		StartingDate:     election.StartingDate,
		FinishingDate:    election.FinishingDate,
		ElectionMode:     election.ElectionMode,
		Voters:           convertVoters(election.Voters),
		PoliticalParties: convertPoliticalParties(election.PoliticalParties),
	}
	return electionGeneric
}

func convertVoters(voters []models2.VoterModel) []modelsGeneric.VoterModel {
	votersGeneric := make([]modelsGeneric.VoterModel, len(voters))
	for _, voter := range voters {
		voterGeneric := modelsGeneric.VoterModel{
			Id:          voter.Id,
			FullName:    voter.Name + " " + voter.LastName,
			BirthDate:   voter.BirthDate,
			Email:       voter.Email,
			Sex:         voter.Sex,
			Phone:       voter.Phone,
			Voted:       voter.Voted,
			OtherFields: map[string]any{"lastname": voter.LastName, "civiccredential": voter.CivicCredential, "department": voter.Department, "circuit": voter.IdCircuit},
		}
		votersGeneric = append(votersGeneric, voterGeneric)
	}
	return votersGeneric
}

func convertPoliticalParties(politicalParties []models2.PoliticalPartyModel) []modelsGeneric.PoliticalPartyModel {
	politicalPartiesGeneric := make([]modelsGeneric.PoliticalPartyModel, len(politicalParties))
	for _, politicalParty := range politicalParties {
		politicalPartyGeneric := modelsGeneric.PoliticalPartyModel{
			Id:          politicalParty.Id,
			Name:        politicalParty.Name,
			Candidates:  convertCandidates(politicalParty.Candidates),
			OtherFields: map[string]any{},
		}
		politicalPartiesGeneric = append(politicalPartiesGeneric, politicalPartyGeneric)
	}
	return politicalPartiesGeneric
}

func convertCandidates(candidates []models2.CandidateModel) []modelsGeneric.CandidateModel {
	candidatesGeneric := make([]modelsGeneric.CandidateModel, len(candidates))
	for _, candidate := range candidates {
		candidateGeneric := modelsGeneric.CandidateModel{
			Id:                 candidate.Id,
			FullName:           candidate.Name + " " + candidate.LastName,
			Sex:                candidate.Sex,
			BirthDate:          candidate.BirthDate,
			IdPoliticalParty:   candidate.IdPoliticalParty,
			NamePoliticalParty: candidate.PoliticalParty,
			OtherFields:        map[string]any{"lastname": candidate.LastName},
		}
		candidatesGeneric = append(candidatesGeneric, candidateGeneric)
	}
	return candidatesGeneric
}
