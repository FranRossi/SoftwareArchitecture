package sms

import (
	"fmt"
	"notification_center/models"
)

func SendCertificateSMS(certificate models.Certificate) {
	fmt.Println("Sending SMS to voter: " + certificate.Fullname + " with phone number: " + certificate.Phone)
	fmt.Println("Certificate: ")
	fmt.Println(certificate)
}
