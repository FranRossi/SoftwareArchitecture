package connections

import (
	conntroller "electoral_api/adapter/uruguayan_election/controller"
	logic2 "electoral_api/adapter/uruguayan_election/logic"
	"electoral_api/adapter/uruguayan_election/repository"
)

func Injection() *conntroller.ElectionController {
	repo := &repository.ElectionRepo{}
	logic := logic2.NewLogicElection(repo)
	return conntroller.NewElectionController(logic)
}
