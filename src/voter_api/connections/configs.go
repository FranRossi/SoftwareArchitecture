package connections

import "os"

func ConfigurationEnvironment() {
	os.Setenv("mq_address", "amqp://guest:guest@localhost:5672/")
	os.Setenv("grpc_address", "localhost:50004")
	os.Setenv("mongo_address", "mongodb://localhost:27017")
}
