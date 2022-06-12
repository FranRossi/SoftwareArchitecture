package logic

import (
	"electoral_service/models"
	"log"
	"voter_api/repository/repository"
)

func FindVoter(idVoter string) (*models.VoterModel, error) {
	user, err := repository.FindVoter(idVoter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, err
}
