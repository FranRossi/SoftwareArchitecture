package dependencies_injection

import (
	"electoral_service/service"
	"electoral_service/service/logic"
	"electoral_service/service/repository"
	"mq"
)

func Injection() *service.ElectionService {
	repo := &repository.ElectionRepo{}
	logic := logic.NewLogicElection(repo)
	service := service.NewElectionService(logic)
	mq.BuildRabbitWorker("amqp://guest:guest@localhost:5672/") // TODO .env
	return service
}
