package models

type Certificate struct {
	IdVoter            string `json:"id_voter"`
	IdElection         string `json:"id_election"`
	TimeVoted          string `json:"time_voted"`
	VoteIdentification string `json:"vote_identification"`
	StartingDate       string `json:"starting_date"`
	FinishingDate      string `json:"finishing_date"`
	ElectionMode       string `json:"election_mode"`
	Fullname           string `json:"fullname"`
	Message            string `json:"message"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}
