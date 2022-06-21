package main

import (
	"external_electoral_api/uruguay_election/connections"
	"external_electoral_api/uruguay_election/repositories"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	repositories.NewUruguayanElection()
	fmt.Println("Connecting")
	connections.Connection()
}
