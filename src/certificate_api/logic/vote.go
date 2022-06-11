package logic

import (
	"fmt"
	domain "voter_api/domain/vote"
	"voter_api/repository/repository"
)

func StoreVote(vote *domain.VoteModel) error {
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	err = repository.RegisterVote(vote.IdVoter)
	return nil
}
