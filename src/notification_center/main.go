package main

import (
	"notification_center/logic"

	"github.com/joho/godotenv"
)

func main() {
	logic.ReceiveActs()
	logic.ReceiveAlert()
	godotenv.Load()

	// fmt.Println("Press Enter to exit")
	// input := bufio.NewScanner(os.Stdin)
	// input.Scan()

	// alert := models.Alert{ //TODO borrar
	// 	IdVoter:    "123",
	// 	IdElection: "df",
	// 	MaxVotes:   123,
	// 	Votes:      333,
	// 	Emails:     []string{"mailDeAlerta1", "mailDeAlerta2"},
	// }
	// logic.NotifyAlertEmails(alert)

}
