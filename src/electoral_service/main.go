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
	err := electoral_service.GetElectionSettings()
	if err != nil {
		return
	}
	mq.GetMQWorker().Close()
}
