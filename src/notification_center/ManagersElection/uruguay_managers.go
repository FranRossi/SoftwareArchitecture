package managersElection

import (
	"fmt"
	"log"
	"notification_center/models"
	"strconv"
)

var emails = []string{"montevideo@intendencia.com", "montevideo@presidencia.com", "colonia@intendencia.com"}

func SendEmails(act models.InitialAct) {
	for _, email := range emails {
		sendEmailTo(email, act)
	}
}

func sendEmailTo(email string, act models.InitialAct) {
	log.Println("Sending email to: " + email)
	fmt.Println()
	log.Println("Comenzó la elección: " + act.StarDate)
	log.Println("La cantidad de votantes que hay habilitados: " + strconv.FormatInt(int64(act.Voters), 10))
	log.Println("El modo de elección es: " + act.Mode)
	for _, politicalParty := range act.PoliticalParties {
		log.Println("La candidaturas de " + politicalParty.Name + " son: ")
		for _, candidate := range politicalParty.Candidates {
			log.Println("El candidato " + candidate.Name + " " + candidate.LastName + "con id " + candidate.Id)
		}
		fmt.Println()
	}
	fmt.Println()
}
