package connections

import (
	mq "message_queue"
	"os"
)

func ConnectionRabbitMQ() {
	mq.BuildRabbitWorker(os.Getenv("mq_address"))
}

func CloseConnectionRabbitMQ() {
	mq.GetMQWorker().Close()
}
