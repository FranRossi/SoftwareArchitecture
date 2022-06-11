package logic

import (
	"log"
	domain "voter_api/domain/user"
	"voter_api/repository/repository"
)

func FindVoter(idVoter string) (*domain.User, error) {
	user, err := repository.FindVoter(idVoter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, err
}

func FindCandidate(idVoter string) (string, error) {
	candidate, err := repository.FindCandidate(idVoter)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return candidate, err
}

func FindElectionMode(idElection string) (string, error) {
	mode, err := repository.FindElectionMode(idElection)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return mode, err
}

func FindElectionTime(idElection string) (string, string, error) {
	startingDate, closingDate, err := repository.FindElectionTime(idElection)
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}
	return startingDate, closingDate, nil
}

func HowManyVotesHasAVoter(idVoter string) int {
	howManyVotesHasAVoter := repository.HowManyVotesHasAVoter(idVoter)
	return howManyVotesHasAVoter
}
