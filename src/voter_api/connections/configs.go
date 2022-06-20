package connections

import (
	"os"
	l "own_logger"
)

func ConfigurationEnvironment() {
	err := os.Setenv("mq_address", "amqp://guest:guest@localhost:5672/")
	err = os.Setenv("grpc_address", ":50004")
	err = os.Setenv("mongo_address", "mongodb://localhost:27017")
	if err != nil {
		l.LogWarning("error getting environmental variables")
	}
}
