package main

import (
	"auth/api"
	"auth/connections"
	"auth/controller"
	"auth/repository"
)

func main() {

	mongoClient := connections.GetInstanceMongoClient()

	repo := repository.NewUsersRepo(mongoClient, "uruguayan_users")
	sessionController := controller.NewSessionController(repo)
	api.ConnectionAPI(sessionController)

	// mongoClient.Disconnect() TODO
}
