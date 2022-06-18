package main

import (
	"bufio"
	"fmt"
	mq "message_queue"
	"os"
	"time"
)

func main() {

	worker := mq.ConnectionRabbit()

	// Listen for messages in queue
	worker.Listen(50, "test-queue", func(message []byte) error {
		// Do something with the body
		time.Sleep(2 * time.Second)
		fmt.Println(string(message))
		return nil
	})

	// Send Messages to queue
	worker.Send("test-queue", []byte("Hello World!"))
	time.Sleep(2 * time.Second)
	worker.Send("test-queue", []byte("Hello World 2!"))
	time.Sleep(4 * time.Second)
	worker.Send("test-queue", []byte("Hello World 3!"))

	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
