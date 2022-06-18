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
	controller := controller.NewSessionController(repo)
	api.ConnectionAPI(controller)

	// mongoClient.Disconnect() TODO
}
