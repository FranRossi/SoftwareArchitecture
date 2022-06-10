package logic

import (
	models2 "electoral_service/adapter/uruguayan_election/models"
	"electoral_service/adapter/uruguayan_election/repository"
	"electoral_service/connections"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type ElectionLogic struct {
	repo *repository.ElectionRepo
}

func NewLogicElection(repo *repository.ElectionRepo) *ElectionLogic {
	return &ElectionLogic{repo: repo}
}

func (logicElection *ElectionLogic) StoreElection(election models2.ElectionModel) error {
	err := logicElection.repo.StoreElectionConfiguration(election)
	if err != nil {
		return fmt.Errorf("election cannot be stored: %w", err)
	}
	err = storeVoters(election.Voters)
	if err != nil {
		return err
	}
	err = storeCandidates(election.PoliticalParties)
	return nil
}

func storeVoters(voters []models2.VoterModel) error {
	err := repository.StoreElectionVoters(voters)
	if err != nil {
		return fmt.Errorf("voters cannot be stored: %w", err)
	}
	return nil
}

func storeCandidates(politicalParties []models2.PoliticalPartyModel) error {
	setPoliticalPartiesNamesToCandidates(politicalParties)
	candidates := politicalParties[0].Candidates
	candidates = append(candidates, politicalParties[1].Candidates...)
	err := repository.StoreCandidates(candidates)
	if err != nil {
		return fmt.Errorf("candidates cannot be stored: %w", err)
	}
	return nil
}

func setPoliticalPartiesNamesToCandidates(politicalParties []models2.PoliticalPartyModel) []models2.PoliticalPartyModel {
	for _, politicalParty := range politicalParties {
		for i, _ := range politicalParty.Candidates {
			politicalParty.Candidates[i].PoliticalParty = politicalParty.Name
		}
	}
	return politicalParties
}

func SetElectionDate(election models2.ElectionModel) {
	startDate, _ := time.Parse(time.RFC3339, election.StartingDate)
	endDate, _ := time.Parse(time.RFC3339, election.FinishingDate)
	setTimer(startDate, startElection(startDate, election.PoliticalParties, len(election.Voters), election.ElectionMode))
	setTimer(endDate, endElection(startDate, endDate, len(election.Voters)))
}

func setTimer(timerDate time.Time, action func()) {
	timer := time.NewTimer(timerDate.Sub(time.Now()))
	done := make(chan bool)
	go func() {
		<-timer.C
		done <- true
	}()
	<-done
	action()
}

func startElection(startDate time.Time, politicalParties []models2.PoliticalPartyModel, voters int, electionMode string) func() {
	return func() {
		fmt.Println("Election started")
		sendInitialAct(startDate, politicalParties, voters, electionMode)
	}
}

func sendInitialAct(startDate time.Time, politicalParties []models2.PoliticalPartyModel, voters int, electionMode string) {
	act := models2.InitialAct{
		StarDate:         startDate.Format(time.RFC3339),
		PoliticalParties: politicalParties,
		Voters:           voters,
		Mode:             electionMode,
	}
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}
	queue := "initial-election-queue"
	connections.ConnectionRabbit(jsonAct, queue)
}

func endElection(startDate, endDate time.Time, voters int) func() {
	return func() {
		resultElection, err := repository.GetVotes()
		if err != nil {
			log.Fatal(err)
		}
		// TODO chequear las validaciones del REQ 2
		fmt.Println("Election finished")
		sendEndingAct(startDate, endDate, voters, resultElection)
	}
}

func sendEndingAct(startDate time.Time, endDate time.Time, voters int, resultElection models2.ResultElection) {
	act := models2.ClosingAct{
		StarDate: startDate.Format(time.RFC3339),
		EndDate:  endDate.Format(time.RFC3339),
		Voters:   voters,
		Result:   resultElection,
	}
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}
	queue := "closing-election-queue"
	connections.ConnectionRabbit(jsonAct, queue)
}

func DropDataBases() {
	repository.DropDataBases()
}
