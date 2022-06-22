package main

import (
	"bufio"
	"consulting_simulator/models"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	l "own_logger"
	"strconv"
	"time"
)

func main() {
	godotenv.Load()
	GenerateConsult()
	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func GenerateConsult() {
	go generateStatistics()
	//generateCertificates()
	//generateVoteCoverage()
	//generateElectionConfiguration()
	//generatePopularVotingTimes()
	//generateElectionResult()
	//generateVote()
}

func generateStatistics() {
	go votesPerRegion()
	go votesPerCircuit()
}

func votesPerRegion() {
	departments := []string{"Montevideo", "Colonia", "Rocha", "Florida", "San Jose", "Soriano", "Salto", "Paysandu", "Treinta y Tres", "Artigas"}
	url := os.Getenv("VOTES_PER_REGION_URL")
	const electionId = "1/"
	url += electionId
	urlBasic := url
	for i := 0; i <= 9; i++ {
		timeFront := time.Now()
		url += departments[i]
		response, err := http.Get(url)
		if err != nil {
			l.LogError(err.Error())
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(response.Body)

		var votesPerRegionCoverage models.VoteRegionCoverage
		er := json.Unmarshal(jsonBytes, &votesPerRegionCoverage)
		if er != nil {
			go l.LogError(err.Error() + " error casting votes coverage per region")
		}
		err = response.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		if votesPerRegionCoverage.Error {
			go l.LogError(votesPerRegionCoverage.Msg + " error on votes per region")
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		message := "Vote coverage per region " + departments[i] + " received correctly on: " + timeReq + " seconds"
		fmt.Println(message)
		l.LogInfo(message)
		url = urlBasic
	}
}

func votesPerCircuit() {
	url := os.Getenv("VOTES_PER_CIRCUIT_URL")
	const electionId = "1/"
	url += electionId
	urlBasic := url
	for i := 0; i <= 9; i++ {
		timeFront := time.Now()
		url += strconv.Itoa(i)
		response, err := http.Get(url)
		if err != nil {
			l.LogError(err.Error())
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(response.Body)

		var votesPerCircuitCoverage models.VoteCircuitCoverage
		er := json.Unmarshal(jsonBytes, &votesPerCircuitCoverage)
		if er != nil {
			go l.LogError(err.Error() + " error casting votes coverage per circuit")
		}
		err = response.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		if votesPerCircuitCoverage.Error {
			go l.LogError(votesPerCircuitCoverage.Msg + " error on votes per circuit")
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		message := "Vote coverage per circuit " + strconv.Itoa(i) + " received correctly on: " + timeReq + " seconds"
		fmt.Println(message)
		l.LogInfo(message)
		url = urlBasic
	}
}
