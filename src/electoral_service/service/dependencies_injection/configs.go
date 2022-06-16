package dependencies_injection

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service"
	"electoral_service/service/logic"
	"electoral_service/service/repository"
)

func Injection() *service.ElectionService {
	repo := &repository.ElectionRepo{}
	adapter := &controller.ElectionController{}
	logicElection := logic.NewLogicElection(repo)
	serviceElection := service.NewElectionService(logicElection, adapter)
	return serviceElection
}
