package logic

import (
	"encoding/json"
	mq "message_queue"
	managersElection "notification_center/ManagersElection/uruguay"
	"notification_center/models"
	l "own_logger"
	"sync"
)

func ReceiveActs() {
	ReceiveAct()
	ReceiveClosingAct()
}

func ReceiveAct() {
	worker := mq.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "initial-election-queue", func(message []byte) error {
		var act models.InitialAct
		er := json.Unmarshal(message, &act)
		if er != nil {
			l.LogError(er.Error())
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
	worker := mq.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "closing-election-queue", func(message []byte) error {
		var act models.ClosingAct
		er := json.Unmarshal(message, &act)
		if er != nil {
			l.LogError(er.Error())
		}
		notifyEmailsClosing(act)
		wg.Done()
		return nil
	})
}

func notifyEmailsClosing(act models.ClosingAct) {
	managersElection.SendClosingEmails(act)
}
