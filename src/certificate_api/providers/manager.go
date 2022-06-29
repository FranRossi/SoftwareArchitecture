package providers

import (
	"certificate_api/models"
	m "electoral_service/models"
	"encoding/json"
	mq "message_queue"
	l "own_logger"
)

type certificateNotification struct {
	IdVoter            string `json:"id_voter"`
	IdElection         string `json:"id_election"`
	TimeVoted          string `json:"time_voted"`
	VoteIdentification string `json:"vote_identification"`
	StartingDate       string `json:"starting_date"`
	FinishingDate      string `json:"finishing_date"`
	ElectionMode       string `json:"election_mode"`
	Fullname           string `json:"fullname"`
	Message            string `json:"message"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}

type AlertCertificates struct {
	VoterId string
}

func SendAlert(voterId string) {
	queue := "certificate-queue-alert"
	alert := AlertCertificates{VoterId: voterId}
	jsonMsg, err := json.Marshal(alert)
	if err != nil {
		l.LogError(err.Error() + "Error while sending to queue " + queue + " json marshalling error")
	}
	mq.GetMQWorker().Send(queue, jsonMsg)
}

func SendSMS(certificate models.CertificateModel, voter m.VoterModel) {
	sendMessageToQueue("certificate-queue-sms", certificate, voter)
}

func SendEmail(certificate models.CertificateModel, voter m.VoterModel) {
	sendMessageToQueue("certificate-queue-email", certificate, voter)
}

func sendMessageToQueue(queueName string, certificate models.CertificateModel, voter m.VoterModel) {
	var msg certificateNotification
	msg.IdVoter = certificate.IdVoter
	msg.IdElection = certificate.IdElection
	msg.TimeVoted = certificate.TimeVoted
	msg.VoteIdentification = certificate.VoteIdentification
	msg.StartingDate = certificate.StartingDate
	msg.FinishingDate = certificate.FinishingDate
	msg.ElectionMode = certificate.ElectionMode
	msg.Fullname = certificate.Fullname
	msg.Phone = voter.Phone
	msg.Email = voter.Email
	msg.Message = certificate.Message

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		l.LogError(err.Error() + "Error while sending to queue " + queueName + " json marshalling error")
	}
	queue := queueName
	mq.GetMQWorker().Send(queue, jsonMsg)
}
