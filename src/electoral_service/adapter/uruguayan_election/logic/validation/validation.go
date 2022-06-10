package validation

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	"fmt"
	p_f "pipes_and_filters"
)

func ValidateInitial(election models2.ElectionModel) error {

	p := p_f.Pipeline{}

	p.LoadFiltersFromYaml("initialValidations.yaml", GetAvailableFilters())

	errors := p.Run(election)

	if len(errors) > 0 {
		// TODO add error code
		return fmt.Errorf("election is not valid: %v", errors)
	}

	return nil
}
