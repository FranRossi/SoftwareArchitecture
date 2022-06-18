package connections

import (
	"fmt"
	mq "message_queue"
	"os"
	l "own_logger"
)

func ConnectionRabbitMQ() {
	_, err := mq.BuildRabbitWorker(os.Getenv("mq_address"))
	if err != nil {
		fmt.Println(err.Error() + "error in building rabbit MQ worker")
		l.LogError(err.Error())
	}
}

func CloseConnectionRabbitMQ() {
	err := mq.GetMQWorker().Close()
	if err != nil {
		fmt.Println(err.Error() + "error in closing connection with rabbit MQ")
		l.LogError(err.Error())
	}
}
