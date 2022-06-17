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
	electionMode, err2 := repository.FindElectionMode(vote.IdElection)
	if err2 != nil {
		return fmt.Errorf("election mode cannot be found: %w", err2)
	}
	howManyTimesVoted := repository.HowManyVotesHasAVoter(vote.IdVoter)
	if electionMode == "multi" && howManyTimesVoted >= 1 {
		errReplacing := updateNewVote(vote)
		if errReplacing != nil {
			return fmt.Errorf("candidate cannot be replaced: %w", errReplacing)
		}
	}
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	err = repository.RegisterVote(vote, electionMode)
	generateElectionSession(vote.IdElection)
	howManyTimesVoted = howManyTimesVoted + 1
	go checkMaxVotesAndSendAlert(howManyTimesVoted, vote)
	return nil
}

func updateNewVote(vote *domain.VoteModel) error {
	err := repository.DeleteOldVote(vote)
	if err != nil {
		return fmt.Errorf("old vote cannot be deleted: %w", err)
	}
	err2 := repository.ReplaceLastCandidateVoted(vote)
	if err2 != nil {
		return fmt.Errorf("vote cannot be updated: %w", err2)
	}
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
	err := repository.StoreVoteInfo(idVoter, idElection, timeFront, timeBack, timePassed, voteIdentification)
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

func SendCertificate(vote *domain.VoteModel, voteIdentification string, timeFront time.Time, err error) {
	timeVoted := timeFront.Format(time.RFC3339)
	certificate := domain.VoteInfo{
		IdVoter:            vote.IdVoter,
		IdElection:         vote.IdElection,
		TimeVoted:          timeVoted,
		VoteIdentification: voteIdentification,
	}
	queue := "voting-certificates"
	if err != nil {
		queue = "voting-certificates-error"
		fmt.Println("error sending certificate: %w", err)
	}
	sendCertificateToMQ(certificate, queue)
}

func sendCertificateToMQ(certificate domain.VoteInfo, queue string) {
	certificateBytes, _ := json.Marshal(certificate)
	mq.GetMQWorker().Send(queue, certificateBytes)
}

func checkMaxVotesAndSendAlert(howManyTimesVoted int, vote *domain.VoteModel) {
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
