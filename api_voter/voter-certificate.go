package api_voter

import (
	"fmt"

	message_queue "239850_221025_219401/workers"
)

func sendCertificate(id string) error {

	worker, err := message_queue.BuildRabbitWorker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	err = worker.Send("certificates-queue", []byte(fmt.Sprintf("Voter with ID %s voted successfully \n", id)))
	if err != nil {
		panic(err)
	}

	return nil
}

func PrintCertificate() error {
	// test function to test recieve message from queue

	worker, err := message_queue.BuildRabbitWorker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	worker.Listen(1, "certificates-queue", func(message []byte) error {
		fmt.Println(string(message))
		return nil
	})

	return nil
}
