package main

import (
	"certificate_api/api"
	"certificate_api/connections"
	"certificate_api/controllers"
	"certificate_api/repositories"
	mq "message_queue"
)

func main() {

	mongoClient := connections.GetInstanceMongoClient()
	mq.BuildRabbitWorker("amqp://guest:guest@localhost:5672/")

	repo := repositories.NewRequestsRepo(mongoClient, "certificates")
	controller := controllers.CertificateRequestsController(repo)
	controllers.ListenerForNewCertificates()
	api.ConnectionAPI(controller)

	mq.GetMQWorker().Close()
	mongoClient.Disconnect()
}
