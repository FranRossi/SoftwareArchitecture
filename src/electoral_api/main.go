package main

import connections "electoral_api/adapter/uruguayan_election/dependencies_injection"

func main() {
	controller := connections.Injection()
	controller.GetElectionSettings()
}
