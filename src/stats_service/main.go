package main

import (
	"bufio"
	"fmt"
	mq "message_queue"
	"os"
	l "own_logger"
	"stats_service/controllers"
)

func main() {
	l.LogInfo("Starting stats application...")
	mq.BuildRabbitWorker("amqp://guest:guest@localhost:5672/")

	controllers.ListenForNewStats()

	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	mq.GetMQWorker().Close()
}
