package own_logger

import (
	"fmt"
	"log"
	"os"
)

func LogError(error string) {
	LogSync(error, "ERROR: ")
}

func LogInfo(info string) {
	LogSync(info, "INFO: ")
}

func LogWarning(warning string) {
	LogSync(warning, "WARNING: ")
}

func LogSync(message, logType string, args ...interface{}) {

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, logType, log.Ldate|log.Ltime|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(message, args...))
	file.Close()
}
