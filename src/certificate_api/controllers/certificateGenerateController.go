package controllers

import (
	"certificate_api/connections"
	"certificate_api/models"
	"certificate_api/providers"
	"certificate_api/repositories"
	"encoding/json"
	"encrypt"
	"fmt"
	mq "message_queue"
	l "own_logger"
)

func ListenerForNewCertificates() {
	mq.GetMQWorker().Listen(50, "voting-certificates", func(message []byte) error {
		var voteInfo models.VoteInfo
		err := json.Unmarshal(message, &voteInfo)
		if err != nil {
			l.LogError(err.Error())
			fmt.Println(err.Error())
			return err // TODO no devolver error, eso solo hace que el rabbit no haga el ack
		}
		mongoClient := connections.GetInstanceMongoClient()
		repo := repositories.NewRequestsRepo(mongoClient, "certificates")
		controller := CertificateRequestsController(repo)
		return controller.GenerateCertificate(voteInfo) // TODO no devolver error, eso solo hace que el rabbit no haga el ack
	})
}

func (controller *CertificateController) GenerateCertificate(voteInfo models.VoteInfo) error {
	var certificate models.CertificateModel
	certificate.IdVoter = voteInfo.IdVoter
	certificate.IdElection = voteInfo.IdElection
	certificate.TimeVoted = voteInfo.TimeVoted
	certificate.VoteIdentification = voteInfo.VoteIdentification

	voter, err := controller.repo.FindVoter(voteInfo.IdVoter)
	if err != nil {
		l.LogError(err.Error())
		fmt.Println(err.Error())
		return fmt.Errorf("voter cannot be found when generating certificate: %w", err)
	}
	election, err := controller.repo.FindElection(voteInfo.IdElection)
	if err != nil {
		l.LogError(err.Error())
		fmt.Println(err.Error())
		return fmt.Errorf("election cannot be found when generating certificate: %w", err)
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
		l.LogError(err.Error())
		fmt.Println(err.Error() + "cannot store certificate")
	}
	return nil
}
