package logic

import (
	"fmt"
	"voter_api/data_access/repository"
	domain "voter_api/domain/vote"
)

func StoreVote(vote *domain.VoteModel) error {
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	return nil
}
