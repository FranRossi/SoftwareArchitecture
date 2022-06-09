package managersElection

import (
	"log"
	"notification_center/models"
	"strconv"
)

var emails = []string{"montevideo@intendencia.com", "montevideo@presidencia.com", "colonia@intendencia.com"}

func SendEmails(act models.Act) {
	for _, email := range emails {
		sendEmailTo(email, act)
	}
}

func sendEmailTo(email string, act models.Act) {
	log.Println("Sending email to: " + email)
	log.Println("Comenzó la elección: " + act.StarDate)
	log.Println("Va a finalizar: " + act.EndDate)
	log.Println("La cantidad de votantes que hay habilitados: " + strconv.FormatInt(int64(act.Voters), 10))
	log.Println("")
}
