package repositories

import (
	"ElectoralService/UruguayanElection"
	"ElectoralService/UruguayanElection/models"
	"fmt"
)

type ElectionRepo struct {
	electionList []models.ElectionModel
}

var electionUruguay models.ElectionModel

func NewUruguayanElection() {
	id, voterAmount := "1", 1000
	var err error
	electionUruguay, err = UruguayanElection.CreateElectionMock(id, voterAmount)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (electionRepo *ElectionRepo) GetElection(id string) (models.ElectionModel, error) {
	//const voterAmount = 10
	//election, err := UruguayanElection.CreateElectionMock(id, voterAmount)
	//if err != nil {
	//	return election, fmt.Errorf("election not found: %s", id)
	//}
	//return election, err
	return electionUruguay, nil
}
