package service

import (
	"os"
	"own_logger"
)

func SetEnvironmentConfig() {
	err := os.Setenv("mq_address", "amqp://guest:guest@localhost:5672/")
	if err != nil {
		own_logger.LogWarning("Error getting environmental variables")
	}
}
