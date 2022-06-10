package logic

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	"fmt"
	p_f "pipes_and_filters"
)

func ValidateInitial(election models2.ElectionModel) error {

	var availableFilters = []p_f.FilterWithName{
		{
			Name:     "validate_election_date",
			Function: FilterEchoData,
		},
		{
			Name:     "validate_party_list",
			Function: FilterEchoData,
		},
		{
			Name:     "validate_candidate_list",
			Function: FilterEchoData,
		},
		{
			Name:     "validate_voter_list",
			Function: FilterEchoData,
		},
		{
			Name:     "validate_party_candidates",
			Function: FilterEchoData,
		},
		{
			Name:     "validate_unique_party_per_candidate",
			Function: FilterEchoData,
		},

		{
			Name:     "validate_election_mode",
			Function: FilterEchoData,
		},
	}

	p := p_f.Pipeline{}

	p.LoadFiltersFromYaml("initialValidations.yaml", availableFilters)

	errors := p.Run(election)

	if len(errors) > 0 {
		// TODO add error code
		return fmt.Errorf("election is not valid: %v", errors)
	}

	return nil
}

func FilterEchoData(data any, params map[string]any) error {
	electionData := data.(models2.ElectionModel)
	fmt.Printf("data Data: %s\n", electionData.Description)
	return fmt.Errorf("Failed")
}
