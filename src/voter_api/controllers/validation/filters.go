package validation

import (
	"fmt"
	p_f "pipes_and_filters"
	domain "voter_api/domain/vote"
	"voter_api/logic"
)

func GetAvailableFilters() map[string]p_f.FilterWithParams {

	availableFilters := map[string]p_f.FilterWithParams{
		"validate_voter":                 FilterValidateVoter,
		"validate_circuit":               FilterValidateCircuit,
		"validate_vote_unique_candidate": FilterValidateUniqueCandidate,
		"validate_candidate":             FilterValidateCandidate,
		"validate_vote_mode":             FilterVoteMode,
	}
	return availableFilters
}

func FilterValidateVoter(data any, params map[string]any) error {
	voter := data.(domain.VoteModel)
	_, err := logic.FindVoter(voter.IdVoter)
	if err != nil {
		return fmt.Errorf("voter is not valid: %v", err)
	}
	return nil
}

func FilterValidateCircuit(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	usr, err := logic.FindVoter(vote.IdVoter)
	if err != nil {
		return fmt.Errorf("voter is not valid: %v", err)
	}
	if usr.IdCircuit != vote.Circuit {
		return fmt.Errorf("voter is not voting on the rigth circuit: %v", err)
	}
	return nil
}

func FilterValidateUniqueCandidate(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	if len(vote.IdCandidate) > 1 {
		return fmt.Errorf("voter can vote only one candidate: %v")
	}
	return nil
}

func FilterValidateCandidate(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	_, err := logic.FindCandidate(vote.IdCandidate)
	if err != nil {
		return fmt.Errorf("candidate is not valid: %v", err)
	}
	return nil
}

func FilterVoteMode(data any, params map[string]any) error {
	//vote := data.(domain.VoteModel)
	//if vote. != "vote" && vote.Mode != "unvote" {
	//	return fmt.Errorf("vote mode is not valid: %v", vote.Mode)
	//}
	return nil
}
