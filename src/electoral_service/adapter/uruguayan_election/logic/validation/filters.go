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
			Function: FilterValidatePoliticalPartyList,
		},
		{
			Name:     "validate_candidate_list",
			Function: FilterValidateCandidateList,
		},
		{
			Name:     "validate_unique_party_per_candidate",
			Function: FilterValidateUniquePartyPerCandidate,
		},
		{
			Name:     "validate_election_mode",
			Function: FilterEchoData,
		},
	}
	return availableFilters
}

func FilterValidatePoliticalPartyList(data any, params map[string]any) error {
	election := data.(models2.ElectionModel)
	if len(election.PoliticalParties) == 0 {
		return fmt.Errorf("politicalPartyList is nil")
	}
	return nil
}

func FilterValidateCandidateList(data any, params map[string]any) error {
	election := data.(models2.ElectionModel)
	for _, party := range election.PoliticalParties {
		if len(party.Candidates) == 0 {
			return fmt.Errorf("candidateList is nil")
		}
	}
	if len(election.Voters) == 0 {
		return fmt.Errorf("voterList is nil")
	}
	return nil
}

func FilterValidateUniquePartyPerCandidate(data any, params map[string]any) error {
	election := data.(models2.ElectionModel)
	candidatesToCheck := make(map[string]bool)
	for _, party := range election.PoliticalParties {
		for _, candidate := range party.Candidates {
			if _, ok := candidatesToCheck[candidate.Id]; ok {
				return fmt.Errorf("candidate %s is in more than one party", candidate.Id)
			} else {
				candidatesToCheck[candidate.Id] = true
			}
		}
	}
	return nil
}
