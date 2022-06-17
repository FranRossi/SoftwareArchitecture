package uruguay

import (
	"fmt"
	"notification_center/models"
	"strconv"
)

func SendAlertEmails(alert models.Alert) {
	for _, email := range alert.Emails {
		sendAlertEmailTo(email, alert)
	}
}

func sendAlertEmailTo(email string, alert models.Alert) {
	fmt.Println("Sending email to: " + email)
	fmt.Println()
	fmt.Println("En la elección: " + alert.IdElection)
	fmt.Println("El votante: " + alert.IdVoter)
	votes := strconv.FormatInt(int64(alert.Votes), 10)
	maxVotes := strconv.FormatInt(int64(alert.MaxVotes), 10)
	fmt.Println("Votó: " + votes + " veces y el máximo es: " + maxVotes)
	fmt.Println()
}
