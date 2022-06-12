package logic

import (
	"fmt"
	"strconv"
	"time"
	"voter_api/controllers/validation"
	domain "voter_api/domain/vote"
	"voter_api/repository/repository"
)

var electionSession = make([]string, 1)

func StoreVote(vote *domain.VoteModel) error {
	validationError := validation.ValidateVote(*vote)
	if validationError != nil {
		return validationError
	}
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	err = repository.RegisterVote(vote.IdVoter)
	generateElectionSession(vote.IdElection)
	return nil
}

func generateElectionSession(idElection string) {
	idElectionInt, _ := strconv.Atoi(idElection)
	electionSession = append(electionSession, strconv.Itoa(idElectionInt))
}

func DeleteVote(vote *domain.VoteModel) error {
	err := repository.DeleteVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be deleted: %w", err)
	}
	return nil
}

func StoreVoteInfo(idVoter, idElection string, timeFrontEnd, timeBackEnd time.Time) (string, error) {
	timeFront := timeFrontEnd.Format(time.RFC3339)
	timeBack := timeBackEnd.Format(time.RFC3339)
	timePassed := timeBackEnd.Sub(timeFrontEnd).String()
	voteIdentification := generateRandomVoteIdentification(idElection)
	err := repository.StoreVoteInfo(idVoter, timeFront, timeBack, timePassed, voteIdentification)
	if err != nil {
		return "", fmt.Errorf("vote info cannot be stored: %w", err)
	}
	return voteIdentification, nil
}

func generateRandomVoteIdentification(idElection string) string {
	idElectionInt, _ := strconv.Atoi(idElection)
	sessionNumber := electionSession[idElectionInt]
	randomNumber := strconv.Itoa(int(time.Now().UnixNano()))
	return sessionNumber + randomNumber
}

type VoteInfo struct {
	IdVoter            string
	IdElection         string
	TimeVoted          string
	VoteIdentification string
}

func SendCertificateSMS(vote *domain.VoteModel, voteIdentification string, timeFront time.Time) {
	timeVoted := timeFront.Format(time.RFC3339)
	certificate := VoteInfo{
		IdVoter:            vote.IdVoter,
		IdElection:         vote.IdElection,
		TimeVoted:          timeVoted,
		VoteIdentification: voteIdentification,
	}
	sendToMQ(certificate)
}

func sendToMQ(certificate VoteInfo) {
	// TODO send to MQ
	fmt.Println("Sending to MQ")
	fmt.Println(certificate)
}
