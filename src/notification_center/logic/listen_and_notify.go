package logic

import (
	"encoding/json"
	mq "message_queue"
	l "own_logger"
)

func listenForMsg[M any](queue_name string, notifyFuncs ...func(act M)) {
	worker := mq.ConnectionRabbit()
	worker.Listen(50, queue_name, func(message []byte) error {
		var act M
		er := json.Unmarshal(message, &act)
		if er != nil {
			l.LogError(er.Error())
		}
		notify(notifyFuncs, act)
		return nil
	})
}

func notify[T any](notifyFuncs []func(act T), act T) {
	for _, notify := range notifyFuncs {
		go notify(act)
	}
}
