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

type InitialAct struct {
	StarDate         string                        `json:"startDate"`
	PoliticalParties []models2.PoliticalPartyModel `json:"politicalParties"`
	Voters           int                           `json:"voters"`
	Mode             string                        `json:"mode"`
}

type ClosingAct struct {
	StarDate   string `json:"startDate"`
	EndDate    string `json:"endDate"`
	Voters     int    `json:"voters"`
	TotalVotes int    `json:"totalVotes"`
	Result     string `json:"result"`
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
	return nil
}

func storeVoters(voters []models2.VoterModel) error {
	err := repository.StoreElectionVoters(voters)
	if err != nil {
		return fmt.Errorf("voters cannot be stored: %w", err)
	}
	return nil
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
	act := InitialAct{
		StarDate:         startDate.Format(time.RFC3339),
		PoliticalParties: politicalParties,
		Voters:           voters,
		Mode:             electionMode,
	}
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}
	connections.ConnectionRabbit(jsonAct)
}

func endElection(startDate, endDate time.Time, voters int) func() {
	return func() {
		fmt.Println("Election finished")
		fmt.Println("Election will finish at: ", endDate)
		//TODO: send closing act
	}
}
