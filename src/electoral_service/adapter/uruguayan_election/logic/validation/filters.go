package validation

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	"fmt"
	p_f "pipes_and_filters"
)

func GetAvailableFilters() []p_f.FilterWithName {

	availableFilters := []p_f.FilterWithName{
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
	return availableFilters
}

func FilterEchoData(data any, params map[string]any) error {
	electionData := data.(models2.ElectionModel)
	fmt.Printf("data Data: %s\n", electionData.Description)
	return fmt.Errorf("Failed")
}
