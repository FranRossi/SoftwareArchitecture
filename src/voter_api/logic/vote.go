package logic

import (
	"fmt"
	"voter_api/controllers/validation"
	domain "voter_api/domain/vote"
	"voter_api/repository/repository"
)

func StoreVote(vote *domain.VoteModel) error {
	validationError := validation.ValidateVote(*vote)
	if validationError != nil {
		return validationError
	}
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	err = repository.RegisterVote(vote.IdVoter)
	return nil
}
