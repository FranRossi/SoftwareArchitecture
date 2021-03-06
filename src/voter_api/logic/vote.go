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
	howManyTimesVoted = howManyTimesVoted + 1
	go checkMaxVotesAndSendAlert(howManyTimesVoted, vote)
	go updateElectionResult(vote, region)
	go repository.RegisterVoteOnCertainGroup(vote.IdElection, voter)
	return nil
}

func updateElectionResult(vote domain.VoteModel, region string) {
	politicalPartyName, errPP := repository.FindPoliticalPartyFromCandidateId(vote.IdCandidate)
	if errPP != nil {
		l.LogError("error getting political party name when storing vote: %w" + errPP.Error())
	}
	err := repository.UpdateElectionResult(vote, region, politicalPartyName)
	if err != nil {
		l.LogError("election result cannot be updated: %w" + err.Error())
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

func StoreVoteInfo(idVoter, idElection, voteIdentification string, timeFrontEnd, timeBackEnd time.Time) error {
	timeFront := timeFrontEnd.Format(time.RFC3339)
	timeBack := timeBackEnd.Format(time.RFC3339)
	timePassed := timeBackEnd.Sub(timeFrontEnd).String()
	err := repository.StoreVoteInfo(idVoter, idElection, timeFront, timeBack, timePassed, voteIdentification)
	if err != nil {
		return fmt.Errorf("vote info cannot be stored: %w", err)
	}
	return nil
}

func GenerateRandomVoteIdentification(idElection string) string {
	randomNumber := strconv.Itoa(int(time.Now().UnixNano()))
	return idElection + randomNumber
}

func SendCertificate(vote domain.VoteModel, voteIdentification string, timeFront time.Time, err error) {
	timeVoted := timeFront.Format(time.RFC3339)
	const message = "vote processed correctly"
	certificate := domain.VoteInfo{
		IdVoter:            vote.IdVoter,
		IdElection:         vote.IdElection,
		TimeVoted:          timeVoted,
		VoteIdentification: voteIdentification,
		Message:            message,
		Error:              false,
	}
	if err != nil {
		certificate.Error = true
		certificate.Message = err.Error() + " ERROR_CODE[VOTE-PROC]"
	}
	queue := "voting-certificates"
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
