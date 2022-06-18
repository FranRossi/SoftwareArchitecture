package uruguayan_election

import (
	"os"
	"own_logger"
)

func ConfigEnvironment() {

	err := os.Setenv("electoral_service_url", "http://localhost:8080/api/v1/election/uruguay/?id=1")
	err = os.Setenv("maxVotes", "1")
	err = os.Setenv("maxCertificate", "10")
	err = os.Setenv("mq_address", "amqp://guest:guest@localhost:5672/")
	if err != nil {
		own_logger.LogWarning("Error getting environmental variables")
	}
}
