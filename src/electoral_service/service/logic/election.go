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
	err := logicElection.storeInitialValues(election)
	if err != nil {
		return err
	}
	return nil
}

func (logicElection *ElectionLogic) storeInitialValues(election *models.ElectionModelEssential) error {
	err := logicElection.repo.StoreElectionConfiguration(election)
	if err != nil {
		return fmt.Errorf("election cannot be stored: %w", err)
	}
	err = storeVoters(election.Voters)
	if err != nil {
		return err
	}
	err = storeCandidates(election.PoliticalParties)
	err = storeInitialResult(election)
	if err != nil {
		return err
	}
	return nil
}

func getAllRegionsFromElection(election *models.ElectionModelEssential) ([]models.Region, error) {
	regionsVoters := make(map[string]int)
	for _, voter := range election.Voters {
		if _, ok := regionsVoters[voter.Region]; !ok {
			regionsVoters[voter.Region] = 0
		}
		regionsVoters[voter.Region]++
	}
	var regionsFromElection []models.Region
	for regionName, amountOfVoterPerRegion := range regionsVoters {
		region := models.Region{
			Name:        regionName,
			Votes:       0,
			TotalVoters: amountOfVoterPerRegion,
		}
		regionsFromElection = append(regionsFromElection, region)
	}
	return regionsFromElection, nil
}

func storeInitialResult(election *models.ElectionModelEssential) error {
	resultElection, err := getElectionResult(election.Id, len(election.Voters))
	regions, err := getAllRegionsFromElection(election)
	resultElection.Regions = regions
	err = repository.StoreElectionResult(resultElection)
	if err != nil {
		return fmt.Errorf("initial result election cannot be stored: %w", err)
	}
	return nil
}

func getElectionResult(electionId string, amountVoters int) (models.ResultElection, error) {
	votesPerCandidates, err := repository.GetEachCandidatesVotes()
	votesPerParties := getVotesPerParties(votesPerCandidates)
	totalVotes, err := repository.GetTotalVotes(electionId)
	if err != nil {
		return models.ResultElection{}, fmt.Errorf("votes election cannot be obtained: %w", err)
	}
	resultElection := models.ResultElection{
		ElectionId:          electionId,
		TotalAmountOfVoters: amountVoters,
		AmountOfVotes:       totalVotes,
		VotesPerParties:     votesPerParties,
		VotesPerCandidates:  votesPerCandidates,
	}
	return resultElection, nil
}

func getVotesPerParties(votesCandidates []models.CandidateEssential) []models.PoliticalPartyEssentials {
	votesPerParties := make(map[string]int, len(votesCandidates))
	for _, candidate := range votesCandidates {
		votesPerParties[candidate.PoliticalParty] += candidate.Votes
	}
	var votesPerPartiesResume []models.PoliticalPartyEssentials
	for key, value := range votesPerParties {
		votesPerPartiesResume = append(votesPerPartiesResume, models.PoliticalPartyEssentials{Name: key, Votes: value})
	}
	return votesPerPartiesResume
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
		resultElection, err := getElectionResult(electionId, voters)
		if err != nil {
			log.Fatal(err)
		}
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
		err = repository.StoreElectionResult(resultElection)
		if err != nil {
			return
		}
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
