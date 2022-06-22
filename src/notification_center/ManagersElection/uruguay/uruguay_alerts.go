package uruguay

import (
	"bufio"
	"fmt"
	"notification_center/models"
	"os"
	l "own_logger"
	"strconv"
)

func SendAlertEmails(alert models.Alert) {

	var emails []string
	emailsFromFile, err := getEmailsFromFile()
	if err != nil || len(emailsFromFile) == 0 {
		emails = alert.Emails
	} else {
		emails = emailsFromFile
	}
	for _, email := range emails {
		l.LogInfo("Sending alert email to: " + email)
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

func getEmailsFromFile() ([]string, error) {
	alertFileName := os.Getenv("ALERT_EMAIL_FILE")
	file, err := os.Open(alertFileName)
	if err != nil {
		l.LogError("could not open file " + alertFileName + "to get emails to send alert")
		return []string{}, err
	}
	defer file.Close()

	var emails []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}
	return emails, scanner.Err()
}
