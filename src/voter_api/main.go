package main

import (
	"voter_api/connections"
	"voter_api/controllers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	connections.ConnectionRabbitMQ()
	grpcServer := connections.ConnectionGRPC()
	jwt := connections.ConnectionJWT()
	voteServer := controllers.RegisterServicesServer(grpcServer, jwt)
	controllers.ActivateChannel(voteServer)
	connections.ServeGRPC(grpcServer)
	defer connections.CloseConnectionRabbitMQ()
}
