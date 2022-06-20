package main

import (
	"bufio"
	"fmt"
	mq "message_queue"
	"os"
	l "own_logger"
)

func main() {
	fmt.Println("Hello, world!")
	l.LogInfo("Starting stats application...")
	mq.BuildRabbitWorker("amqp://guest:guest@localhost:5672/")

	//repo := repositories.NewRequestsRepo(mongoClient, "certificates")

	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	mq.GetMQWorker().Close()
}
