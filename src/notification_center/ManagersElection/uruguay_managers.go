package managersElection

import (
	"fmt"
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
	fmt.Println("Sending email to: " + email)
	fmt.Println()
	fmt.Println("Comenzó la elección: " + act.StartDate)
	fmt.Println("La cantidad de votantes que hay habilitados: " + strconv.FormatInt(int64(act.Voters), 10))
	fmt.Println("El modo de elección es: " + act.Mode)
	for _, politicalParty := range act.PoliticalParties {
		fmt.Println("La candidaturas del " + politicalParty.Name + " son: ")
		for _, candidate := range politicalParty.Candidates {
			fmt.Println("El candidato " + candidate.Name + " " + candidate.LastName + "con id " + candidate.Id)
		}
		fmt.Println()
	}
	fmt.Println()
}
