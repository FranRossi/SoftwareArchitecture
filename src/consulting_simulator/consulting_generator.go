package main

import (
	"bufio"
	"bytes"
	"consulting_simulator/models"
	"consulting_simulator/repository"
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
	client := GetInstanceMongoClient()
	repo := repository.NewRequestsRepo(client, "certificates")
	certificates, err := repo.FindAllCertificateForElection("1")
	if err != nil {
		l.LogError(err.Error())
	}
	rolConsulter := "Consulter"
	rolElectoral := "Electoral"
	rolConsultingAgent := "ConsultingAgents"
	tokens := RegisterAndLoginUser(rolConsulter, rolElectoral, rolConsultingAgent)
	GenerateConsult(certificates, tokens)
	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func GenerateConsult(certificates []models.CertificateResponseModel, tokens []string) {
	go generateStatistics()
	go generateCertificates(certificates)
	go generateElectionConfiguration(tokens[0])
	go generatePopularVotingTimes()
	go generateElectionResult(tokens[1])
	go generateVote(tokens[1])
}

func generateVote(tokenElectoral string) {
	url := os.Getenv("SPECIFIC_VOTE_URL")
	const electionId = "1"
	url += electionId + "/"
	urlBasic := url
	for i := 0; i <= 80; i++ {
		timeFront := time.Now()
		voterId := strconv.FormatInt(int64(i+10000000), 10)
		url += voterId
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", tokenElectoral)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			l.LogError("There are no election's votes with for that voter at the moment, retry later")
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		var vote models.VoteJson
		er := json.Unmarshal(jsonBytes, &vote)
		if er != nil {
			go l.LogError(err.Error() + " error casting vote ")
		}
		err = resp.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if vote.Error {
			message := "Vote could not be found for voter: " + voterId + " for election: " + electionId + " response on: " + timeReq + " seconds"
			l.LogError(message)
			fmt.Println(message)
		} else {
			message := "Election result for election with id " + electionId + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}
		url = urlBasic
	}

}

func generateElectionResult(tokenElectoral string) {
	url := os.Getenv("ELECTION_RESULT_URL")
	const electionId = "1"
	url += electionId
	for i := 0; i <= 30; i++ {
		timeFront := time.Now()
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", tokenElectoral)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			l.LogError("There are no election's result at the moment, retry later")
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		var electionResult models.ElectionResultJson
		er := json.Unmarshal(jsonBytes, &electionResult)
		if er != nil {
			go l.LogError(err.Error() + " error casting election result")
		}
		err = resp.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if electionResult.Error {
			go l.LogError(electionResult.Msg + " error on election's result")
			fmt.Println("There is not result for election: " + electionId)
		} else {
			message := "Election result for election with id " + electionId + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}
	}
}

func generatePopularVotingTimes() {
	url := os.Getenv("VOTER_PER_HOUR_URL")
	const electionId = "1/"
	url += electionId
	for i := 0; i <= 15; i++ {
		timeFront := time.Now()
		response, err := http.Get(url)
		if err != nil {
			l.LogError("There are not votes per hours available, retry later")
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(response.Body)

		var votesPerHours models.VoterPerHoursJson
		er := json.Unmarshal(jsonBytes, &votesPerHours)
		if er != nil {
			go l.LogError(err.Error() + " error casting votes per hours")
		}
		err = response.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if votesPerHours.Error {
			go l.LogError(votesPerHours.Msg + " error on votes per hours")
			fmt.Println("There are not voter per hours for election: " + electionId)
		} else {
			message := "Vote per hours received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}

	}
}

func RegisterAndLoginUser(rolConsulter, rolElectoral, rolConsultingAgent string) []string {
	var tokens []string
	userConsulter := models.UserRegister{
		Id:       "10000000",
		Role:     rolConsulter,
		Password: "1234",
	}
	userElectoral := models.UserRegister{
		Id:       "10000001",
		Role:     rolElectoral,
		Password: "1234",
	}
	userConsultingAgent := models.UserRegister{
		Id:       "10000002",
		Role:     rolConsultingAgent,
		Password: "1234",
	}
	Register(userConsulter)
	Register(userConsultingAgent)
	Register(userElectoral)
	tokenConsulter := Login(models.Login{Id: userConsulter.Id, Password: userConsulter.Password})
	tokenElectoral := Login(models.Login{Id: userElectoral.Id, Password: userElectoral.Password})
	tokenConsultingAgent := Login(models.Login{Id: userConsultingAgent.Id, Password: userConsultingAgent.Password})
	tokens = append(tokens, tokenConsulter, tokenElectoral, tokenConsultingAgent)
	return tokens
}

func Login(credentials models.Login) string {
	urlLogin := os.Getenv("LOGIN_USER_URL")
	loginUserConsulter := map[string]string{"id": credentials.Id, "password": credentials.Password}
	json_data, err := json.Marshal(loginUserConsulter)
	respTokenJson, err := http.Post(urlLogin, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		l.LogError(err.Error())
		return ""
	}
	jsonBytes, err := ioutil.ReadAll(respTokenJson.Body)
	var token models.Token
	er := json.Unmarshal(jsonBytes, &token)
	if er != nil {
		go l.LogError(err.Error() + " error casting token")
	}
	err = respTokenJson.Body.Close()
	if err != nil {
		go l.LogError(err.Error())
	}
	return token.Token
}

func Register(user models.UserRegister) {
	urlRegister := os.Getenv("REGISTER_USER_URL")
	jsonDataConsulter, err := json.Marshal(user)
	_, err = http.Post(urlRegister, "application/json", bytes.NewBuffer(jsonDataConsulter))
	if err != nil {
		l.LogError(err.Error())
		return
	}
}

func generateElectionConfiguration(tokenConsulter string) {
	url := os.Getenv("ELECTION_CONFIG_URL")

	const electionId = "1"
	url += electionId
	for i := 0; i <= 30; i++ {
		timeFront := time.Now()
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", tokenConsulter)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			l.LogError("There are no election available with that id, retry later")
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		var electionConfig models.ElectionConfigJson
		er := json.Unmarshal(jsonBytes, &electionConfig)
		if er != nil {
			go l.LogError(err.Error() + " error casting election config")
		}
		err = resp.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if electionConfig.Error {
			go l.LogError(electionConfig.Msg + " error on election config")
			fmt.Println("There are not config for election: " + electionId)
		} else {
			message := "Election config with id " + electionId + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}
	}
}

func generateCertificates(certificates []models.CertificateResponseModel) {
	url := os.Getenv("CERTIFICATES_URL")
	for i := 0; i <= 80; i++ {
		timeFront := time.Now()
		certificate := models.CertificateRequestModel{
			VoterId:            certificates[i].IdVoter,
			VoteIdentification: certificates[i].VoteIdentification,
		}
		jsonData, err := json.Marshal(certificate)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			l.LogError(err.Error())
			return
		}
		timeBack := time.Now()
		jsonBytes, err := ioutil.ReadAll(response.Body)

		var certificateRequest models.CertificateJson
		er := json.Unmarshal(jsonBytes, &certificateRequest)
		if er != nil {
			go l.LogError(err.Error() + " error casting certificates")
		}
		err = response.Body.Close()
		if err != nil {
			go l.LogError(err.Error())
		}
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if certificateRequest.Error {
			go l.LogError(certificateRequest.Msg + " error on certificate")
			fmt.Println("There are not certificates for voter: " + certificate.VoterId)
		} else {
			message := "Certificates for  " + certificates[i].IdVoter + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}

	}
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
			l.LogError("There are not votes coverage per regions available, retry later")
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
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if votesPerRegionCoverage.Error {
			go l.LogError(votesPerRegionCoverage.Msg + " error on votes per region")
			fmt.Println("There are not votes for that circuit")
		} else {
			message := "Vote coverage per region " + departments[i] + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}
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
		timesSub := timeBack.Sub(timeFront).Seconds()
		timeReq := fmt.Sprintf("%f", timesSub)
		if votesPerCircuitCoverage.Error {
			go l.LogError(votesPerCircuitCoverage.Msg + " error on votes per circuit")
			fmt.Println("There are not votes for that circuit")
		} else {
			message := "Vote coverage per circuit " + strconv.Itoa(i) + " received correctly on: " + timeReq + " seconds"
			fmt.Println(message)
			l.LogInfo(message)
		}
		url = urlBasic
	}
}
