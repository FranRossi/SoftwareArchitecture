package repositories

import (
	"ElectoralService/UruguayanElection/models"
	"fmt"
)

type ElectionRepo struct {
	electionList []models.Election
}

func (electionRepo *ElectionRepo) SendElectionSettings(election models.Election) error {
	electionRepo.electionList = append(electionRepo.electionList, election)
	return nil
}

func (electionRepo *ElectionRepo) GetElection(id string) (models.Election, error) {
	for _, election := range electionRepo.electionList {
		if election.Id == id {
			return election, nil
		}
	}
	return models.Election{}, fmt.Errorf("election id not found: %s", id)
}
