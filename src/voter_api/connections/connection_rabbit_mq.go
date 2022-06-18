package connections

import (
	mq "message_queue"
	"os"
	l "own_logger"
)

func ConnectionRabbitMQ() {
	_, err := mq.BuildRabbitWorker(os.Getenv("mq_address"))
	if err != nil {
		l.LogError(err.Error())
	}
}

func CloseConnectionRabbitMQ() {
	err := mq.GetMQWorker().Close()
	if err != nil {
		l.LogError(err.Error())
	}
}
