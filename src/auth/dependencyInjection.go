package main

import (
	"auth/api"
	"auth/connections"
	"auth/controller"
	"auth/jwt"
	"auth/repository"
	"io/ioutil"
	"os"
	l "own_logger"
	"time"
)

func InjectDependencies() {
	mongoClient := connections.GetInstanceMongoClient()

	repo := repository.NewUsersRepo(mongoClient, os.Getenv("DATABASE_NAME"))
	manager := createJwtManager()
	sessionController := controller.NewSessionController(repo, manager)
	api.ConnectionAPI(sessionController)
	defer connections.CloseMongoClient()
}

func createJwtManager() *jwt.Manager {
	duration := 30 * time.Minute

	privateKey, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		l.LogError(err.Error())
	}
	publicKey, err := ioutil.ReadFile("./public.rsa")
	if err != nil {
		l.LogError(err.Error())
	}
	manager := jwt.NewJWTManager(privateKey, publicKey, duration)
	return manager
}
