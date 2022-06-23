package email

import (
	"fmt"
	"notification_center/models"
	l "own_logger"
	"strconv"
)

func SendInitialActsEmails(act models.InitialAct) {
	for _, email := range act.Emails {
		l.LogInfo("Sending email with initial act to: " + email)
		sendInitialEmailTo(email, act)
	}
}

func sendInitialEmailTo(email string, act models.InitialAct) {
	fmt.Println("Sending email to: " + email)
	fmt.Println()
	fmt.Println("Comenzó la elección: " + act.ElectionId + " " + act.StartDate)
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

func SendClosingEmails(act models.ClosingAct) {
	for _, email := range act.Emails {
		l.LogInfo("Sending email with closing act to: " + email)
		sendClosingEmailTo(email, act)
	}
}

func sendClosingEmailTo(email string, act models.ClosingAct) {
	fmt.Println("Sending email to: " + email)
	fmt.Println()
	fmt.Println("Comenzó la elección: " + act.ElectionId + " " + act.StarDate)
	fmt.Println("Finalizó la elección: " + act.EndDate)
	fmt.Println("La cantidad de votantes que ha votado: " + strconv.FormatInt(int64(act.Voters), 10))
	fmt.Println("Los resultados de la elección son: ")
	for _, party := range act.Result.VotesPerParties {
		fmt.Println("La cantidad de votos de " + party.Name + " es: " + strconv.FormatInt(int64(party.Votes), 10))
	}
	fmt.Println()
	for _, candidate := range act.Result.VotesPerCandidates {
		fmt.Println("El candidato " + candidate.Name + " " + "con id " + candidate.Id + " ha obtenido " + strconv.FormatInt(int64(candidate.Votes), 10) + " votos")
	}
	fmt.Println()
}
