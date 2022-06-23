package service

import (
	"electoral_service/models"
	"electoral_service/service/logic"
	"fmt"
	l "own_logger"
)

type Adapter interface {
	GetElectionSettings() *models.ElectionModelEssential
}

type ElectionService struct {
	adapter       Adapter
	electionLogic *logic.ElectionLogic
}

func NewElectionService(logic *logic.ElectionLogic, adapter Adapter) *ElectionService {
	return &ElectionService{
		electionLogic: logic,
		adapter:       adapter,
	}
}

func (service *ElectionService) GetElectionSettings() {
	election := service.adapter.GetElectionSettings()
	err := service.electionLogic.StoreElection(election)
	if err != nil {
		l.LogError(err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Election stored successfully")
	l.LogInfo("Election stored successfully")

	logic.SetElectionDate(election)
}

func (service *ElectionService) DropDataBases() {
	service.electionLogic.DropDataBases()
}
