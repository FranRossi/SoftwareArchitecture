package main

import (
	"voter_api/connections"
	"voter_api/controllers"
)

func main() {
	connections.ConfigurationEnvironment()
	connections.ConnectionRabbitMQ()
	grpcServer := connections.ConnectionGRPC()
	jwt := connections.ConnectionJWT()
	voteServer := controllers.RegisterServicesServer(grpcServer, jwt)
	controllers.ActivateChannel(voteServer)
	connections.ServeGRPC(grpcServer)
	defer connections.CloseConnectionRabbitMQ()
}
