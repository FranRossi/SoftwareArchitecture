package logic

import (
	"encoding/json"
	"log"
	managersElection "notification_center/ManagersElection"
	connections "notification_center/connection"
	"notification_center/models"
	"sync"
)

func RecieveAct() {
	worker := connections.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "election-settings-queue", func(message []byte) error {
		var act models.Act
		er := json.Unmarshal(message, &act)
		if er != nil {
			log.Fatal(er)
		}
		notifyEmails(act)
		wg.Done()
		return nil
	})

}

func notifyEmails(act models.Act) {
	managersElection.SendEmails(act)
}
