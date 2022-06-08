package logic

import (
	"encoding/json"
	"fmt"
	"log"
	connections "notification_center/connection"
	"sync"
)

type Act struct {
	starDate string
	endDate  string
	voters   int
}

func RecieveAct() {
	worker := connections.ConnectionRabbit()
	wg := sync.WaitGroup{}
	worker.Listen(50, "election-settings-queue", func(message []byte) error {

		var act Act
		er := json.Unmarshal(message, &act)
		if er != nil {
			log.Fatal(er)
		}
		fmt.Println(act.starDate, act.endDate, act.voters)
		wg.Done()
		return nil
	})

}
