package repositories

import "ElectoralService/UruguayanElection/models"

type ElectionRepo struct {
	electionList []models.Election
}

func (electionRepo *ElectionRepo) SendElectionSettings(election models.Election) error {
	electionRepo.electionList = append(electionRepo.electionList, election)
	return nil
}
