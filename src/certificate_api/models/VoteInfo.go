package models

type VoteInfo struct {
	IdVoter            string
	IdElection         string
	TimeVoted          string
	VoteIdentification string
	Message            string
	Error              bool
}
