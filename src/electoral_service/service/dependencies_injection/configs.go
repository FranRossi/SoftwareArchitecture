package dependencies_injection

import (
	"electoral_service/adapter/uruguayan_election/controller"
	"electoral_service/service"
	"electoral_service/service/logic"
	"electoral_service/service/repository"
	"message_queue"
	"os"
)

func Injection() *service.ElectionService {
	repo := &repository.ElectionRepo{}
	adapter := &controller.ElectionController{}
	logic := logic.NewLogicElection(repo)
	service := service.NewElectionService(logic, adapter)
	mq.BuildRabbitWorker(os.Getenv("mq_address"))
	return service
}
