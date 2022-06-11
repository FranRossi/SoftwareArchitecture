package main

import (
	"voter_api/connections"
	"voter_api/controllers"
)

func main() {
	grpcServer := connections.ConnectionGRPC()
	jwt := connections.ConnectionJWT()
	controllers.RegisterServicesServer(grpcServer, jwt)
	connections.ServeGRPC(grpcServer)
}
