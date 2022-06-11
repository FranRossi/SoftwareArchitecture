package service

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service/logic"
	"fmt"
	"log"
	"time"
)

const url = "http://localhost:8080/api/v1/election/uruguay/?id=1"

type ElectionService struct {
	adapter       *controller.ElectionController // TODO change to interface, and use dependency injection, to inject the adapter
	electionLogic *logic.ElectionLogic
}

func NewElectionService(logic *logic.ElectionLogic) *ElectionService {
	return &ElectionService{electionLogic: logic}
}

func (service *ElectionService) GetElectionSettings() error {
	election := service.adapter.GetElectionSettings()
	// TODO fix external service and remove following lines:
	election.StartingDate = time.Now().Add(time.Second * 30).Format(time.RFC3339)
	election.FinishingDate = time.Now().Add(time.Minute * 1).Format(time.RFC3339)
	err := service.electionLogic.StoreElection(election)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Election stored successfully")
	logic.SetElectionDate(election)
	return nil
}

func (service *ElectionService) DropDataBases() {
	service.electionLogic.DropDataBases()
}
