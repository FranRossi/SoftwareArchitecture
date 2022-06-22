package main

import (
	"consulting_simulator/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	l "own_logger"
)

func main() {
	GenerateConsult()
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
	votesPerRegion()
}

func votesPerRegion() {
	url := os.Getenv("VOTES_PER_REGION_URL")
	const electionId = "1"
	for i := 0; i <= 10; i++ {
		url += electionId + "/" + string(i)
		response, err := http.Get(url)
		if err != nil {
			l.LogError(err.Error())
		}
		jsonBytes, err := ioutil.ReadAll(response.Body)
		defer func() {
			err := response.Body.Close()
			if err != nil {
				go l.LogError(err.Error())
			}
		}()
		var votesPerRegionCoverage models.VoteRegionCoverage
		er := json.Unmarshal(jsonBytes, &votesPerRegionCoverage)
		if er != nil {
			go l.LogError(err.Error() + " error casting votes coverage per region")
		}
		if votesPerRegionCoverage.Error {
			go l.LogError(votesPerRegionCoverage.Msg + " error on votes per region")
		}
	}
}
