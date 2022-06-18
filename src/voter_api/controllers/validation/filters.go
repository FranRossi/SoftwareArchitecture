package validation

import (
	"fmt"
	p_f "pipes_and_filters"
	"time"
	"voter_api/domain"
	"voter_api/repository"
)

func GetAvailableFilters() map[string]p_f.FilterWithParams {

	availableFilters := map[string]p_f.FilterWithParams{
		"validate_voter":                 FilterValidateVoter,
		"validate_circuit":               FilterValidateCircuit,
		"validate_vote_unique_candidate": FilterValidateUniqueCandidate,
		"validate_candidate":             FilterValidateCandidate,
		"validate_vote_mode":             FilterVoteMode,
		"validate_voting_time":           FilterValidateVotingTime,
	}
	return availableFilters
}

func FilterValidateVoter(data any, params map[string]any) error {
	voter := data.(domain.VoteModel)
	_, err := repository.FindVoter(voter.IdVoter)
	if err != nil {
		return fmt.Errorf("voter is not valid: %v", err)
	}
	return nil
}

func FilterValidateCircuit(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	usr, err := repository.FindVoter(vote.IdVoter)
	if err != nil {
		return fmt.Errorf("voter is not valid: %v", err)
	}
	if usr.OtherFields["circuit"].(string) != vote.Circuit {
		return fmt.Errorf("voter is not voting on the rigth circuit")
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
	_, err := repository.FindCandidate(vote.IdCandidate)
	if err != nil {
		return fmt.Errorf("candidate is not valid: %v", err)
	}
	return nil
}

func FilterVoteMode(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	modeExpected := "unico"
	mode, err := repository.FindElectionMode(vote.IdElection)
	if err == nil && mode == modeExpected {
		howManyVotesHasAVoter := repository.HowManyVotesHasAVoter(vote.IdVoter)
		if howManyVotesHasAVoter > 0 {
			return fmt.Errorf("voter has already voted")
		}
	}
	if err != nil {
		return fmt.Errorf("election mode is not valid: %v", err)
	}
	return nil
}

func FilterValidateVotingTime(data any, params map[string]any) error {
	vote := data.(domain.VoteModel)
	startingDate, closingDate, err := repository.FindElectionTime(vote.IdElection)
	if err != nil {
		return fmt.Errorf("election does not exist: %v", err)
	}
	startingDateAsDate, err := time.Parse(time.RFC3339, startingDate)
	closingDateAsDate, err := time.Parse(time.RFC3339, closingDate)
	if err != nil {
		return fmt.Errorf("election date are misconfigured: %v", err)
	}
	now := time.Now()
	if now.Before(startingDateAsDate) {
		return fmt.Errorf("election has not yet started")
	}
	if now.After(closingDateAsDate) {
		return fmt.Errorf("election is over")
	}
	return nil
}
