package service

import (
	"own_logger"

	"github.com/joho/godotenv"
)

func SetEnvironmentConfig() {
	err := godotenv.Load()
	if err != nil {
		own_logger.LogWarning("Error getting environmental variables")
	}
}
