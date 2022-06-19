package main

import (
	"auth/jwt"
	"consulting_api/api"
	"consulting_api/connections"
	"consulting_api/controllers"
	"consulting_api/repositories"
	"io/ioutil"
	l "own_logger"
	"time"
)

func InjectDependencies() {
	mongoClient := connections.GetInstanceMongoClient()
	manager := createJwtManager()

	repo := repositories.NewRequestsRepo(mongoClient, "uruguay_votes")
	controller := controllers.NewConsultingController(repo, manager)
	repoElection := repositories.NewElectionRepo(mongoClient, "uruguay_election")
	electionController := controllers.NewConsultingElectionConfigController(repoElection, manager)
	api.ConnectionAPI(controller, electionController)
}

func createJwtManager() *jwt.Manager {
	duration := 30 * time.Minute

	privateKey, err := ioutil.ReadFile("./jwt_key.rsa")
	if err != nil {
		l.LogError(err.Error())
	}
	publicKey, err := ioutil.ReadFile("./jwt_public_key.rsa")
	if err != nil {
		l.LogError(err.Error())
	}
	manager := jwt.NewJWTManager(privateKey, publicKey, duration)
	return manager
}
