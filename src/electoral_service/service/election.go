package service

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service/logic"
	"fmt"
	"log"
)

type ElectionService struct {
	adapter       *controller.ElectionController // TODO change to interface, and use dependency injection, to inject the adapter
	electionLogic *logic.ElectionLogic
}

func NewElectionService(logic *logic.ElectionLogic, adapter *controller.ElectionController) *ElectionService {
	return &ElectionService{
		electionLogic: logic,
		adapter:       adapter,
	}
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

func (service *ElectionService) DropDataBases() {
	service.electionLogic.DropDataBases()
}
