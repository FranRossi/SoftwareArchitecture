package logic

import (
	"encoding/json"
	mq "message_queue"
	l "own_logger"
)

func listenForMsg[M any](queue_name string, notifyFuncs ...func(msg M)) {
	worker := mq.ConnectionRabbit()
	worker.Listen(50, queue_name, func(message []byte) error {
		var model M
		er := json.Unmarshal(message, &model)
		if er != nil {
			l.LogError(er.Error())
		}
		notify(notifyFuncs, model)
		return nil
	})
}

func notify[T any](notifyFuncs []func(msg T), msg T) {
	for _, notify := range notifyFuncs {
		go notify(msg)
	}
}
