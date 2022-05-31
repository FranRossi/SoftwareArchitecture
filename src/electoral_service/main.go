package main

import (
	"electoral_service/UruguayanElection/connections"
	"electoral_service/UruguayanElection/repositories"
	"fmt"
)

func main() {
	repositories.NewUruguayanElection()
	fmt.Println("Connecting")
	connections.Connection()
}
