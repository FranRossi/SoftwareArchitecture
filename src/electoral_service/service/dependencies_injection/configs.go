package dependencies_injection

import (
	"electoral_service/service"
	"electoral_service/service/logic"
	"electoral_service/service/repository"
)

func Injection() *service.ElectionService {
	repo := &repository.ElectionRepo{}
	logic := logic.NewLogicElection(repo)
	service := service.NewElectionService(logic)
	return service
}
