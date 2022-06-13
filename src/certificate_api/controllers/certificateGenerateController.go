package controllers

import (
	"certificate_api/models"
	"certificate_api/repositories"
	"encoding/json"
	"log"
	mq "message_queue"
)

//funcion que recibe el modelo de info de la votacion y genera el certificado

func ListenerForNewCertificates() {
	mq.GetMQWorker().Listen(50, "voting-certificates", func(message []byte) error {
		var voteInfo models.VoteInfo
		err := json.Unmarshal(message, &voteInfo)
		if err != nil {
			log.Fatal(err)
			return err
		}

		return GenerateCertificate(voteInfo)
	})
}

func GenerateCertificate(voteInfo models.VoteInfo) error {
	log.Println(voteInfo)
	var certificate models.CertificateModel
	certificate.IdVoter = voteInfo.IdVoter
	certificate.IdElection = voteInfo.IdElection
	certificate.TimeVoted = voteInfo.TimeVoted
	certificate.VoteIdentification = voteInfo.VoteIdentification

	fullName, _ := repositories.FindVoterFullName(voteInfo.IdVoter)
	certificate.Fullname = fullName

	// certificate.StartingDate =
	// certificate.FinishingDate =
	// certificate.ElectionMode =
	// generate certificate and send sms
	return nil
}
