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
	controllers.RegisterServicesServer(grpcServer, jwt)
	connections.ServeGRPC(grpcServer)
	defer connections.CloseConnectionRabbitMQ()
}
