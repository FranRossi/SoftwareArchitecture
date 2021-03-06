package validation

import (
	"fmt"
	"own_logger"
	p_f "pipes_and_filters"
	"voter_api/domain"
)

func ValidateVote(vote domain.VoteModel) error {
	p := p_f.Pipeline{}
	errLoadingYaml := p.LoadFiltersFromYaml("voteValidations.yaml", GetAvailableFilters())
	if errLoadingYaml != nil {
		return errLoadingYaml
	}

	errors := p.Run(vote)
	if len(errors) > 0 {
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
