package email

import (
	"fmt"
	"notification_center/models"
)

func SendCertificateEmail(certificate models.Certificate) {
	fmt.Println("Sending Email to voter: " + certificate.Fullname + " with email: " + certificate.Email)
	fmt.Println("Certificate: ")
	fmt.Println(certificate)
}
