package logic

import (
	"encoding/json"
	mq "message_queue"
	managersElection "notification_center/ManagersElection/uruguay"
	"notification_center/models"
	l "own_logger"
	"sync"
)

func ReceiveAlert() {
	worker := mq.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "alert-queue", func(message []byte) error {
		var alert models.Alert
		er := json.Unmarshal(message, &alert)
		if er != nil {
			l.LogError(er.Error())
		}
		NotifyAlertEmails(alert)
		wg.Done()
		return nil
	})
}

func NotifyAlertEmails(alert models.Alert) {
	managersElection.SendAlertEmails(alert)
}
