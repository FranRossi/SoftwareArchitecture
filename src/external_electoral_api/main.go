package main

import (
	"external_electoral_api/uruguay_election/connections"
	"external_electoral_api/uruguay_election/repositories"
	"fmt"
)

func main() {
	repositories.NewUruguayanElection()
	fmt.Println("Connecting")
	connections.Connection()
}
