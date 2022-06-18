package main

import (
	dependencyinjection "electoral_service/service/dependencies_injection"
	mq "message_queue"
	"os"
)

func main() {
	mq.BuildRabbitWorker(os.Getenv("mq_address"))
	electoral_service := dependencyinjection.Injection()
	electoral_service.DropDataBases()
	electoral_service.GetElectionSettings()
	mq.GetMQWorker().Close()
}
