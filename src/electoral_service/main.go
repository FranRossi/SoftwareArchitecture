package main

import (
	dependencyinjection "electoral_service/service/dependencies_injection"
)

func main() {
	controller := dependencyinjection.Injection()
	controller.DropDataBases()
	controller.GetElectionSettings()
}
