package main

import (
	"auth/api"
	"auth/connections"
	"auth/controller"
	"auth/repository"
)

func main() {

	mongoClient := connections.GetInstanceMongoClient()

	repo := repositories.NewUsersRepo(mongoClient, "uruguayan_users")
	controller := controllers.NewSessionController(repo)
	api.ConnectionAPI(controller)

	// mongoClient.Disconnect() TODO
}
