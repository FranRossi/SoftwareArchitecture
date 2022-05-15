package api_voter

import (
	"fmt"
)

func sendCertificate(id string) error {
	fmt.Printf("Voter with ID %s voted successfully \n", id)
	return nil
}
