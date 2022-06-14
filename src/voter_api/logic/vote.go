package logic

import (
	"encoding/json"
	"fmt"
	mq "message_queue"
	"strconv"
	"time"
	"voter_api/controllers/validation"
	"voter_api/domain"
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
	howManyTimesVoted := repository.HowManyVotesHasAVoter(vote.IdVoter)
	go checkMaxVotes(howManyTimesVoted, vote)
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

func SendCertificate(vote *domain.VoteModel, voteIdentification string, timeFront time.Time) {
	timeVoted := timeFront.Format(time.RFC3339)
	certificate := VoteInfo{
		IdVoter:            vote.IdVoter,
		IdElection:         vote.IdElection,
		TimeVoted:          timeVoted,
		VoteIdentification: voteIdentification,
	}
	sendCertificateToMQ(certificate)
}

func sendCertificateToMQ(certificate VoteInfo) {
	queue := "voting-certificates"
	certificateBytes, _ := json.Marshal(certificate)
	mq.GetMQWorker().Send(queue, certificateBytes)
}

func checkMaxVotes(howManyTimesVoted int, vote *domain.VoteModel) {
	maxVotes, _, err := repository.GetMaximumValuesBeforeAlert(vote.IdElection)
	if err != nil {
		fmt.Println(err)
	}
	if howManyTimesVoted >= maxVotes {
		sendAlertToMQ(vote.IdVoter, vote.IdElection, howManyTimesVoted, maxVotes)
	}
}

func sendAlertToMQ(idVoter, idElection string, howManyTimesVoted int, maxVotes int) {
	queue := "alert-queue"
	alert := domain.Alert{
		IdVoter:    idVoter,
		IdElection: idElection,
		MaxVotes:   maxVotes,
		Votes:      howManyTimesVoted,
	}
	alertBytes, _ := json.Marshal(alert)
	mq.GetMQWorker().Send(queue, alertBytes)
}
