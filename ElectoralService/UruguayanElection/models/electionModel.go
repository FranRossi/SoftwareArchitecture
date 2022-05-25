package models

type Election struct {
	Id            string `json:"id"`
	Description   string `json:"description"`
	StartingDate  string `json:"startingDate"`
	FinishingDate string `json:"finishingDate"`
	ElectionMode  string `json:"electionMode"`
}
