package main

import dependencyinjection "electoral_service/adapter/uruguayan_election/dependencies_injection"

func main() {
	controller := dependencyinjection.Injection()
	controller.DropDataBases()
	controller.GetElectionSettings()
}
