package dependencies_injection

import (
	"electoral_service/service"
	logic2 "electoral_service/service/logic"
	"electoral_service/service/repository"
)

func Injection() *service.ElectionService {
	repo := &repository.ElectionRepo{}
	logic := logic2.NewLogicElection(repo)
	return service.NewElectionService{}
}
