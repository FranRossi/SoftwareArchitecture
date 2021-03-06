package main

import (
	"electoral_service/connections"
	"electoral_service/service"
	dependencyinjection "electoral_service/service/dependencies_injection"
	mq "message_queue"
	"os"
)

func main() {
	service.SetEnvironmentConfig()
	mq.BuildRabbitWorker(os.Getenv("MQ_HOST"))

	electoral_service := dependencyinjection.Injection()

	electoral_service.DropDataBases()
	electoral_service.GetElectionSettings()
	mq.GetMQWorker().Close()
	defer connections.CloseMongoClient()
}
