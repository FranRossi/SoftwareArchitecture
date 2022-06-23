package main

import (
	"bufio"
	"fmt"
	mq "message_queue"
	"os"
	l "own_logger"
	"stats_service/controllers"
	"stats_service/repository"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	repository.DropDataBases()
	l.LogInfo("Starting stats application...")
	mq.BuildRabbitWorker(os.Getenv("MQ_HOST"))

	controllers.ListenForNewStats()

	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	mq.GetMQWorker().Close()
}
