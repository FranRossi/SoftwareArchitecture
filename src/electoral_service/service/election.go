package service

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service/logic"
	"fmt"
	"log"
)

const url = "http://localhost:8080/api/v1/election/uruguay/?id=1"

type ElectionService struct {
	adapter       *controller.ElectionController
	electionLogic *logic.ElectionLogic
}

func NewElectionService(logic *logic.ElectionLogic) *ElectionService {
	return &ElectionService{electionLogic: logic}
}

func (service *ElectionService) GetElectionSettings() error {
	election := service.adapter.GetElectionSettings()
	err := service.electionLogic.StoreElection(election)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Election stored successfully")
	logic.SetElectionDate(election)
	return nil
}
