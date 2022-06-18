package validation

import (
	"fmt"
	"own_logger"
	p_f "pipes_and_filters"
	"voter_api/domain"
)

func ValidateVote(vote domain.VoteModel) error {
	p := p_f.Pipeline{}
	p.LoadFiltersFromYaml("voteValidations.yaml", GetAvailableFilters())
	errors := p.Run(vote)
	if len(errors) > 0 {
		// TODO add error code
		LogValidationErrors(errors)
		return fmt.Errorf("vote is not valid: %v", errors)
	}
	return nil
}

func LogValidationErrors(errors []error) {
	for _, er := range errors {
		own_logger.LogError(er.Error())
	}
}
