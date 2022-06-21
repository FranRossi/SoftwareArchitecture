package validation

import (
	"electoral_service/models"
	"fmt"
	"own_logger"
	p_f "pipes_and_filters"
)

func ValidateInitial(election models.ElectionModelEssential) error {
	p := p_f.Pipeline{}
	errLoadingYaml := p.LoadFiltersFromYaml("initialValidations.yaml", GetAvailableFilters())
	if errLoadingYaml != nil {
		return errLoadingYaml
	}
	errors := p.Run(election)
	if len(errors) > 0 {
		LogValidationErrors(errors)
		//TODO add error code
		return fmt.Errorf("election is not valid: %v", errors)
	}
	return nil
}

func ValidateEndAct(act models.ClosingAct) error {
	p := p_f.Pipeline{}
	errLoadingYaml := p.LoadFiltersFromYaml("endValidations.yaml", GetAvailableFilters())
	if errLoadingYaml != nil {
		return errLoadingYaml
	}
	errors := p.Run(act)
	if len(errors) > 0 {
		LogValidationErrors(errors)
		//TODO add error code
		return fmt.Errorf("end final: %v", errors)
	}
	return nil
}

func LogValidationErrors(errors []error) {
	for _, er := range errors {
		own_logger.LogError(er.Error())
	}
}
