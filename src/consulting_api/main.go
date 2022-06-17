package main

import (
	"consulting_api/api"
	"consulting_api/connections"
	"consulting_api/controllers"
	"consulting_api/repositories"
)

func main() {
	mongoClient := connections.GetInstanceMongoClient()
	repo := repositories.NewRequestsRepo(mongoClient, "uruguay_votes")
	controller := controllers.NewConsultingController(repo)
	repoElection := repositories.NewElectionRepo(mongoClient, "uruguay_election")
	electionController := controllers.NewConsultingElectionConfigController(repoElection)
	api.ConnectionAPI(controller, electionController)
	// mongoClient.Disconnect() TODO
}
