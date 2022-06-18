package logic

import (
	"encoding/json"
	"fmt"
	mq "message_queue"
	l "own_logger"
	"strconv"
	"time"
	"voter_api/controllers/validation"
	"voter_api/domain"
	"voter_api/repository"
)

var electionSession = make([]string, 1)

func StoreVote(vote domain.VoteModel) error {
	validationError := validation.ValidateVote(vote)
	if validationError != nil {
		return validationError
	}
	electionMode, err2 := repository.FindElectionMode(vote.IdElection)
	if err2 != nil {
		return fmt.Errorf("election mode cannot be found when storing vote: %w", err2)
	}
	voter, err3 := repository.FindVoter(vote.IdVoter)
	region := voter.Region
	if err3 != nil {
		return fmt.Errorf("error getting voter region when storing vote: %w", err3)
	}
	howManyTimesVoted := repository.HowManyVotesHasAVoter(vote.IdVoter)
	if electionMode == "multi" && howManyTimesVoted >= 1 {
		errReplacing := updateNewVote(vote, region)
		if errReplacing != nil {
			return fmt.Errorf("candidate cannot be replaced: %w", errReplacing)
		}
	}
	err := repository.StoreVote(vote)
	if err != nil {
		return fmt.Errorf("vote cannot be stored: %w", err)
	}
	err4 := repository.RegisterVote(vote, electionMode)
	if err4 != nil {
		return fmt.Errorf("vote cannot be registered: %w", err4)
	}
	generateElectionSession(vote.IdElection)
	howManyTimesVoted = howManyTimesVoted + 1
	politicalParty, err5 := repository.FindPoliticalPartyFromCandidateId(vote.IdCandidate)
	if err5 != nil {
		return fmt.Errorf("error getting political party name when storing vote: %w", err5)
	}
	go checkMaxVotesAndSendAlert(howManyTimesVoted, vote)
	go updateElectionResult(vote, politicalParty, region)
	return nil
}

func updateElectionResult(vote domain.VoteModel, politicalPartyName, region string) {
	err := repository.UpdateElectionResult(vote, region, politicalPartyName)
	if err != nil {
		fmt.Errorf("election result cannot be updated: %w", err)
	}
}

func updateNewVote(vote domain.VoteModel, region string) error {
	err := repository.DeleteOldVote(vote, region)
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

func DeleteVote(vote domain.VoteModel) error {
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

func SendCertificate(vote domain.VoteModel, voteIdentification string, timeFront time.Time, err error) {
	timeVoted := timeFront.Format(time.RFC3339)
	certificate := domain.VoteInfo{
		IdVoter:            vote.IdVoter,
		IdElection:         vote.IdElection,
		TimeVoted:          timeVoted,
		VoteIdentification: voteIdentification,
	}
	queue := "voting-certificates"
	if err != nil {
		l.LogError(err.Error())
		queue = "voting-certificates-error"
		fmt.Println("error sending certificate: %w", err)
	}
	sendCertificateToMQ(certificate, queue)
}

func sendCertificateToMQ(certificate domain.VoteInfo, queue string) {
	certificateBytes, _ := json.Marshal(certificate)
	err := mq.GetMQWorker().Send(queue, certificateBytes)
	if err != nil {
		l.LogError(err.Error())
	}
}

func checkMaxVotesAndSendAlert(howManyTimesVoted int, vote domain.VoteModel) {
	maxVotes, emails, err := repository.GetMaxVotesAndEmailsBeforeAlert(vote.IdElection)
	if err != nil {
		fmt.Println(err)
	}
	if howManyTimesVoted >= maxVotes {
		sendAlertToMQ(vote, howManyTimesVoted, maxVotes, emails)
	}
}

func sendAlertToMQ(vote domain.VoteModel, howManyTimesVoted, maxVotes int, emails []string) {
	queue := "alert-queue"
	alert := domain.Alert{
		IdVoter:    vote.IdVoter,
		IdElection: vote.IdElection,
		MaxVotes:   maxVotes,
		Votes:      howManyTimesVoted,
		Emails:     emails,
	}
	alertBytes, _ := json.Marshal(alert)
	err := mq.GetMQWorker().Send(queue, alertBytes)
	if err != nil {
		l.LogError(err.Error())
	}
}
