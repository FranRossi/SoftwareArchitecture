package models

type Alert struct {
	IdVoter    string
	IdElection string
	MaxVotes   int
	Votes      int
	Emails     []string
}
