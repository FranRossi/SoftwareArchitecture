package main

import (
	"ElectoralService/UruguayanElection/connections"
	"ElectoralService/UruguayanElection/repositories"
	"fmt"
)

func main() {
	repositories.NewUruguayanElection()
	fmt.Println("Connecting")
	connections.Connection()
}
