package logic

import (
	"notification_center/models"
	"notification_center/providers/email"
	"notification_center/providers/sms"
)

func StartReceivingMsgs() {

	// ADD PROVIDER'S FUNCTIONS AS PARAMETER ON THE CORRESPONDING RECEIVER
	receiveInitialAct(email.SendInitialActsEmails)
	receiveClosingAct(email.SendClosingEmails)
	receiveVotesAlert(email.SendVotesAlertEmails)
	receiveCertificateGenerated(sms.SendCertificateSMS)
	receiveCertificateRequested(email.SendCertificateEmail)
	receiveCertificateAlerts(email.SendCertificatesAlertEmails)
}

func receiveInitialAct(notifyFuncs ...func(act models.InitialAct)) {
	listenForMsg("initial-election-queue", notifyFuncs...)
}

func receiveClosingAct(notifyFuncs ...func(act models.ClosingAct)) {
	listenForMsg("initial-election-queue", notifyFuncs...)
}

func receiveVotesAlert(notifyFuncs ...func(act models.AlertVotes)) {
	listenForMsg("alert-queue", notifyFuncs...)
}

func receiveCertificateGenerated(notifyFuncs ...func(act models.Certificate)) {
	listenForMsg("certificate-queue-sms", notifyFuncs...)
}

func receiveCertificateRequested(notifyFuncs ...func(act models.Certificate)) {
	listenForMsg("certificate-queue-email", notifyFuncs...)
}

func receiveCertificateAlerts(notifyFuncs ...func(alert models.AlertCertificates)) {
	listenForMsg("certificate-queue-alert", notifyFuncs...)
}
