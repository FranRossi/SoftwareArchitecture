package connections

import (
	conntroller "electoral_service/adapter/uruguayan_election/controller"
	logic2 "electoral_service/adapter/uruguayan_election/logic"
	"electoral_service/adapter/uruguayan_election/repository"
)

func Injection() *conntroller.ElectionController {
	repo := &repository.ElectionRepo{}
	logic := logic2.NewLogicElection(repo)
	return conntroller.NewElectionController(logic)
}
