package logic

import (
	"electoral_service/models"
	"electoral_service/service/logic/validation"
	"electoral_service/service/repository"
	"encoding/json"
	"fmt"
	"log"
	mq "message_queue"
	"time"
)

type ElectionLogic struct {
	repo *repository.ElectionRepo
}

func NewLogicElection(repo *repository.ElectionRepo) *ElectionLogic {
	return &ElectionLogic{repo: repo}
}

func (logicElection *ElectionLogic) StoreElection(election *models.ElectionModelEssential) error {

	validationError := validation.ValidateInitial(*election)
	if validationError != nil {
		return validationError
	}

	err := logicElection.repo.StoreElectionConfiguration(election)
	if err != nil {
		return fmt.Errorf("election cannot be stored: %w", err)
	}
	err = storeVoters(election.Voters)
	if err != nil {
		return err
	}
	err = storeCandidates(election.PoliticalParties)
	return nil
}

func storeVoters(voters []models.VoterModel) error {
	err := repository.StoreElectionVoters(voters)
	if err != nil {
		return fmt.Errorf("voters cannot be stored: %w", err)
	}
	return nil
}

func storeCandidates(politicalParties []models.PoliticalPartyModel) error {
	setPoliticalPartiesNamesToCandidates(politicalParties)
	candidates := politicalParties[0].Candidates
	candidates = append(candidates, politicalParties[1].Candidates...)
	err := repository.StoreCandidates(candidates)
	if err != nil {
		return fmt.Errorf("candidates cannot be stored: %w", err)
	}
	return nil
}

func setPoliticalPartiesNamesToCandidates(politicalParties []models.PoliticalPartyModel) []models.PoliticalPartyModel {
	for _, politicalParty := range politicalParties {
		for i := range politicalParty.Candidates {
			politicalParty.Candidates[i].NamePoliticalParty = politicalParty.Name
		}
	}
	return politicalParties
}

func SetElectionDate(election *models.ElectionModelEssential) {
	startDate, _ := time.Parse(time.RFC3339, election.StartingDate)
	endDate, _ := time.Parse(time.RFC3339, election.FinishingDate)
	setTimer(startDate, startElection(startDate, election.PoliticalParties, len(election.Voters), election.ElectionMode, election.Id))
	setTimer(endDate, endElection(startDate, endDate, len(election.Voters), election.Id))
}

func setTimer(timerDate time.Time, action func()) {
	timer := time.NewTimer(time.Until(timerDate))
	done := make(chan bool)
	go func() {
		<-timer.C
		done <- true
	}()
	<-done
	action()
}

func startElection(startDate time.Time, politicalParties []models.PoliticalPartyModel, voters int, electionMode, electionId string) func() {
	return func() {
		fmt.Println("Election started")
		sendInitialAct(startDate, politicalParties, voters, electionMode, electionId)
	}
}

func sendInitialAct(startDate time.Time, politicalParties []models.PoliticalPartyModel, voters int, electionMode, electionId string) {
	act := models.InitialAct{
		StarDate:         startDate.Format(time.RFC3339),
		PoliticalParties: politicalParties,
		Voters:           voters,
		Mode:             electionMode,
		ElectionId:       electionId,
	}
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}

	queue := "initial-election-queue"
	mq.GetMQWorker().Send(queue, jsonAct)
}

func endElection(startDate, endDate time.Time, voters int, electionId string) func() {
	return func() {
		resultElection, err := repository.GetVotes(electionId)
		if err != nil {
			log.Fatal(err)
		}
		resultElection.TotalAmountOfVoters = voters
		resultElection.ElectionId = electionId
		act := models.ClosingAct{
			StarDate: startDate.Format(time.RFC3339),
			EndDate:  endDate.Format(time.RFC3339),
			Result:   resultElection,
		}

		validationError := validation.ValidateEndAct(act)
		if validationError != nil {
			log.Fatal(validationError)
		}
		fmt.Println("Election finished")
		go sendEndingAct(act)
		repository.StoreElectionResult(resultElection)
	}
}

func sendEndingAct(act models.ClosingAct) {
	jsonAct, err := json.Marshal(act)
	if err != nil {
		log.Fatal(err)
	}
	queue := "closing-election-queue"
	mq.GetMQWorker().Send(queue, jsonAct)
}

func (logicElection *ElectionLogic) DropDataBases() {
	repository.DropDataBases()
}
