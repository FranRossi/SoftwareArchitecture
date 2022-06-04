package logic

import (
	models2 "electoral_api/adapter/uruguayan_election/models"
	"electoral_api/adapter/uruguayan_election/repository"
	"fmt"
)

func StoreElection(election models2.ElectionModel) error {
	err := repository.StoreElectionConfiguration(election)
	if err != nil {
		return fmt.Errorf("election cannot be stored: %w", err)
	}
	err = storeVoters(election.Voters)
	if err != nil {
		return err
	}

	return nil
}

func storeVoters(voters []models2.VoterModel) error {
	var votersInterface []interface{}

	for _, v := range voters {
		votersInterface = append(votersInterface, v)
	}
	err := repository.StoreElectionVoters(votersInterface)
	if err != nil {
		return fmt.Errorf("voters cannot be stored: %w", err)
	}
	return nil
}
