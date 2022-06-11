package validation

import (
	"electoral_service/models"
	"fmt"
	p_f "pipes_and_filters"
)

func ValidateInitial(election models.ElectionModelEssential) error {
	p := p_f.Pipeline{}
	p.LoadFiltersFromYaml("initialValidations.yaml", GetAvailableFilters())
	errors := p.Run(election)
	if len(errors) > 0 {
		// TODO add error code
		return fmt.Errorf("election is not valid: %v", errors)
	}
	return nil
}
