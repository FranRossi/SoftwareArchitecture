package uruguayan_election

import "os"

func ConfigEnvironment() {
	os.Setenv("electoral_service_url", "http://localhost:8080/api/v1/election/uruguay/?id=1")

	os.Setenv("maxVotes", "5")
	os.Setenv("maxCertificate", "10")
	os.Setenv("mq_address", "amqp://guest:guest@localhost:5672/")
}
