package logic

import (
	"encoding/json"
	"log"
	managersElection "notification_center/ManagersElection/uruguay"
	connections "notification_center/connection"
	"notification_center/models"
	"sync"
)

func RecieveActs() {
	ReceiveAct()
	ReceiveClosingAct()
}

func ReceiveAct() {
	worker := connections.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "initial-election-queue", func(message []byte) error {
		var act models.InitialAct
		er := json.Unmarshal(message, &act)
		if er != nil {
			log.Fatal(er)
		}
		notifyEmails(act)
		wg.Done()
		return nil
	})
}

func notifyEmails(act models.InitialAct) {
	managersElection.SendInitialActsEmails(act)
}

func ReceiveClosingAct() {
	worker := connections.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "closing-election-queue", func(message []byte) error {
		var act models.ClosingAct
		er := json.Unmarshal(message, &act)
		if er != nil {
			log.Fatal(er)
		}
		notifyEmailsClosing(act)
		wg.Done()
		return nil
	})
}

func notifyEmailsClosing(act models.ClosingAct) {
	managersElection.SendClosingEmails(act)
}
