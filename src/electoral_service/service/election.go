package service

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service/logic"
	"fmt"
	"own_logger"
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

func (service *ElectionService) GetElectionSettings() {
	election := service.adapter.GetElectionSettings()
	err := service.electionLogic.StoreElection(election)
	if err != nil {
		own_logger.LogError(err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Election stored successfully")
	logic.SetElectionDate(election)
}

func (service *ElectionService) DropDataBases() {
	service.electionLogic.DropDataBases()
}
