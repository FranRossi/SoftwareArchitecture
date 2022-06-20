package controllers

import (
	mq "message_queue"
)

func ListenerForStatsForTotal() {
	mq.GetMQWorker().Listen(50, "stats-total", func(message []byte) error {

		return nil
	})
}
