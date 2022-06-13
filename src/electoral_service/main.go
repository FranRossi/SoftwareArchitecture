package main

import (
	dependencyinjection "electoral_service/service/dependencies_injection"
	"mq"
)

func main() {
	electoral_service := dependencyinjection.Injection()
	electoral_service.DropDataBases()
	electoral_service.GetElectionSettings()
	mq.GetMQWorker().Close()
}
