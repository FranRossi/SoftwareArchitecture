package main

import connections "electoral_service/adapter/uruguayan_election/dependencies_injection"

func main() {
	controller := connections.Injection()
	controller.GetElectionSettings()
}
