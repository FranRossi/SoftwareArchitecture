package logic

import (
	"encoding/json"
	"log"
	managersElection "notification_center/ManagersElection/uruguay"
	connections "notification_center/connection"
	"notification_center/models"
	"sync"
)

func ReceiveAlert() {
	worker := connections.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "alert-queue", func(message []byte) error {
		var alert models.Alert
		er := json.Unmarshal(message, &alert)
		if er != nil {
			log.Fatal(er)
		}
		notifyAlertEmails(alert)
		wg.Done()
		return nil
	})
}

func notifyAlertEmails(alert models.Alert) {
	managersElection.SendAlertEmails(alert)
}
