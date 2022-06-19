package providers

import (
	"certificate_api/models"
	m "electoral_service/models"
	"fmt"
)

func SendSMS(certificate models.CertificateModel, voter m.VoterModel) {
	fmt.Println("Sending SMS to voter: " + voter.FullName + " with phone: " + voter.Phone)
	fmt.Println("Certificate: ")
	fmt.Println(certificate)
}

func SendEmail(certificate models.CertificateModel, voter m.VoterModel) {
	fmt.Println("Sending Email to voter: " + voter.FullName + " with email: " + voter.Email)
	fmt.Println("Certificate: ")
	fmt.Println(certificate)
}
