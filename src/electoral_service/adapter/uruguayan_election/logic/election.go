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

type Act struct {
	starDate string
	endDate  string
	voters   int
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

func SetElectionDate(startingDate string, finishingDate string, voters int) {
	startDate, _ := time.Parse(time.RFC3339, startingDate)
	endDate, _ := time.Parse(time.RFC3339, finishingDate)
	setTimer(startDate, startElection(startDate, endDate, voters))
	setTimer(endDate, endElection(endDate))
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

func startElection(startDate time.Time, endDate time.Time, voters int) func() {
	return func() {
		fmt.Println("Election started")
		sendInitialAct(startDate, endDate, voters)
	}
}

func sendInitialAct(startDate time.Time, endDate time.Time, voters int) {
	act := Act{
		starDate: startDate.Format(time.RFC3339),
		endDate:  endDate.Format(time.RFC3339),
		voters:   voters,
	}
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}
	println(jsonAct)
	connections.ConnectionRabbit(jsonAct)
}

func endElection(endDate time.Time) func() {
	return func() {
		fmt.Println("Election finished")
		fmt.Println("Election will finish at: ", endDate)
	}
}
