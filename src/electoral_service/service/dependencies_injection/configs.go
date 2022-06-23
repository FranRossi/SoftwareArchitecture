package dependencies_injection

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service"
	"electoral_service/service/logic"
	"electoral_service/service/repository"
	"os"
)

func Injection() *service.ElectionService {

	var availableAdapters = map[string]service.Adapter{
		"uruguay": &controller.UruguayAdapter{},
	}
	adapter := availableAdapters[os.Getenv("ADAPTER")]

	repo := &repository.ElectionRepo{}
	logicElection := logic.NewLogicElection(repo)
	serviceElection := service.NewElectionService(logicElection, adapter)
	return serviceElection
}
