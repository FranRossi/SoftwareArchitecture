package repository

import (
	"electoral_service/models"
	"encoding/json"
	mq "message_queue"
	l "own_logger"
)

type VoterStats struct {
	ElectionId string
	BirthDate  string
	Region     string
	Circuit    string
	Sex        string
}

func RegisterVoteOnCertainGroup(electionId string, voter *models.VoterModel) {
	fullVoter, errFinding := FindVoter(voter.Id)
	if errFinding != nil {
		return
	}
	var voterStats VoterStats
	voterStats.BirthDate = fullVoter.BirthDate
	voterStats.Circuit = fullVoter.OtherFields["circuit"].(string)
	voterStats.Region = fullVoter.Region
	voterStats.Sex = fullVoter.Sex
	voterStats.ElectionId = electionId

	jsonStats, errs := json.Marshal(voterStats)
	if errs != nil {
		l.LogError("error sending voter stats to queue:" + errs.Error())
	}
	mq.GetMQWorker().Send("stats-actual", jsonStats)
}
