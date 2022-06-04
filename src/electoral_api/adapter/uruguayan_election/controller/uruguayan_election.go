package controller

import (
	models2 "electoral_api/adapter/uruguayan_election/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "http://localhost:8080/api/v1/election/uruguay/?id=1"

func GetElectionSettings() {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var electionSettings models2.ElectionJson
	er := json.Unmarshal(jsonBytes, &electionSettings)
	if er != nil {
		log.Fatal(er)
	}
	fmt.Println(electionSettings.Election)
}
