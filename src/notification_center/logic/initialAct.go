package logic

import (
	"encoding/json"
	"fmt"
	"log"
	connections "notification_center/connection"
	"sync"
)

type Act struct {
	StarDate string `json:"startDate"`
	EndDate  string `json:"endDate"`
	Voters   int    `json:"voters"`
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
		fmt.Printf("%+v", act)
		wg.Done()
		return nil
	})

}
