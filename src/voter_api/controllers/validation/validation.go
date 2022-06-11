package validation

import (
	"fmt"
	p_f "pipes_and_filters"
	domain "voter_api/domain/vote"
)

func ValidateVote(vote domain.VoteModel) error {
	p := p_f.Pipeline{}
	p.LoadFiltersFromYaml("voteValidations.yaml", GetAvailableFilters())
	errors := p.Run(vote)
	if len(errors) > 0 {
		// TODO add error code
		return fmt.Errorf("vote is not valid: %v", errors)
	}
	return nil
}
