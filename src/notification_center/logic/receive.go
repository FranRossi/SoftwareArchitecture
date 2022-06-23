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
	receiveAlert(email.SendAlertEmails)
	receiveCertificateGenerated(sms.SendCertificateSMS)
	receiveCertificateRequested(email.SendCertificateEmail)
}

func receiveInitialAct(notifyFuncs ...func(act models.InitialAct)) {
	listenForMsg("initial-election-queue", notifyFuncs...)
}

func receiveClosingAct(notifyFuncs ...func(act models.ClosingAct)) {
	listenForMsg("initial-election-queue", notifyFuncs...)
}

func receiveAlert(notifyFuncs ...func(act models.Alert)) {
	listenForMsg("alert-queue", notifyFuncs...)
}

func receiveCertificateGenerated(notifyFuncs ...func(act models.Certificate)) {
	listenForMsg("certificate-queue-sms", notifyFuncs...)
}

func receiveCertificateRequested(notifyFuncs ...func(act models.Certificate)) {
	listenForMsg("certificate-queue-email", notifyFuncs...)
}
