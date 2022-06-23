package main

import (
	"certificate_api/api"
	"certificate_api/connections"
	"certificate_api/controllers"
	"certificate_api/repositories"
	mq "message_queue"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	mongoClient := connections.GetInstanceMongoClient()
	mq.BuildRabbitWorker(os.Getenv("MQ_HOST"))

	repo := repositories.NewRequestsRepo(mongoClient, os.Getenv("CERTIFICATES_DB"))
	repo.DropDataBases()
	controller := controllers.CertificateRequestsController(repo)
	controllers.ListenerForNewCertificates()
	api.ConnectionAPI(controller)

	mq.GetMQWorker().Close()
	connections.CloseMongoClient()
}
