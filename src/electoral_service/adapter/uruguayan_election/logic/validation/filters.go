package validation

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	"fmt"
	p_f "pipes_and_filters"
	"time"
)

func GetAvailableFilters() []p_f.FilterWithName {

	availableFilters := []p_f.FilterWithName{
		{
			Name:     "validate_election_date",
			Function: FilterValidateDate,
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
			Function: FilterElectionMode,
		},
	}
	return availableFilters
}

func FilterEchoData(data any, params map[string]any) error {
	electionData := data.(models2.ElectionModel)
	fmt.Printf("data Data: %s\n", electionData.Description)
	return fmt.Errorf("Failed")
}

func FilterValidateDate(data any, params map[string]any) error {
	election := data.(models2.ElectionModel)
	StartingDate, _ := time.Parse(time.RFC3339, election.StartingDate)
	EndDate, _ := time.Parse(time.RFC3339, election.FinishingDate)

	if StartingDate.After(EndDate) || StartingDate.Equal(EndDate) {
		return fmt.Errorf("starting date is after end date")
	}
	return nil
}

func FilterElectionMode(data any, params map[string]any) error {
	election := data.(models2.ElectionModel)
	if election.ElectionMode != "unico" && election.ElectionMode != "multi" {
		return fmt.Errorf("election mode is not valid")
	}
	return nil
}
