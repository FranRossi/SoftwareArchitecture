package controllers

import (
	"certificate_api/connections"
	"certificate_api/models"
	"certificate_api/providers"
	"certificate_api/repositories"
	"encoding/json"
	"encrypt"
	mq "message_queue"
	l "own_logger"
)

func ListenerForNewCertificates() {
	mq.GetMQWorker().Listen(50, "voting-certificates", func(message []byte) error {
		var voteInfo models.VoteInfo
		err := json.Unmarshal(message, &voteInfo)
		if err != nil {
			l.LogError(err.Error())
			return nil
		}
		mongoClient := connections.GetInstanceMongoClient()
		repo := repositories.NewRequestsRepo(mongoClient, "certificates")
		controller := CertificateRequestsController(repo)
		controller.GenerateCertificate(voteInfo)
		return nil
	})
}

func (controller *CertificateController) GenerateCertificate(voteInfo models.VoteInfo) {
	var certificate models.CertificateModel
	certificate.IdVoter = voteInfo.IdVoter
	certificate.IdElection = voteInfo.IdElection
	certificate.TimeVoted = voteInfo.TimeVoted
	certificate.VoteIdentification = voteInfo.VoteIdentification
	certificate.Message = voteInfo.Message
	voter, err := controller.repo.FindVoter(voteInfo.IdVoter)
	if err != nil {
		l.LogError("voter cannot be found when generating certificate: " + err.Error())
	}
	election, err := controller.repo.FindElection(voteInfo.IdElection)
	if err != nil {
		l.LogError("election cannot be found when generating certificate: " + err.Error())
	}
	encrypt.DecryptVoter(&voter)
	certificate.Fullname = voter.FullName
	l.LogInfo("Generating certificate for voter: " + voter.FullName)
	certificate.StartingDate = election.StartingDate
	certificate.FinishingDate = election.FinishingDate
	certificate.ElectionMode = election.ElectionMode
	go providers.SendSMS(certificate, voter)
	err = controller.repo.StoreCertificate(certificate)
	if err != nil {
		l.LogError(err.Error() + "cannot store certificate")
	}
}
