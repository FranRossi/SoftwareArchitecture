package controller

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	modelsGeneric "electoral_service/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	l "own_logger"
)

type UruguayAdapter struct {
}

func (controller *UruguayAdapter) GetElectionSettings() *modelsGeneric.ElectionModelEssential {
	url := os.Getenv("electoral_service_url")
	response, err := http.Get(url)
	if err != nil {
		l.LogError(err.Error())
	}
	jsonBytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
	}()
	var electionSettings models2.ElectionJson
	er := json.Unmarshal(jsonBytes, &electionSettings)
	if er != nil {
		go l.LogError(err.Error() + " error casting election settings")
	}
	if electionSettings.Error {
		go l.LogError(electionSettings.Msg + " error on election settings from external service")
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
		OtherFields:      convertMaximumVotesAndCertificateAndEmails(election.Emails),
	}
	return electionGeneric
}

func convertMaximumVotesAndCertificateAndEmails(emails []string) map[string]any {
	otherFields := map[string]any{
		"maxVotes":       os.Getenv("maxVotes"),
		"maxCertificate": os.Getenv("maxCertificate"),
		"emails":         emails,
	}
	return otherFields
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
