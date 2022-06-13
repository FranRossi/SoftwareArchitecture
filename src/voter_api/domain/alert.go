package domain

type Alert struct {
	IdVoter    string
	IdElection string
	MaxVotes   int
	Votes      int
}
