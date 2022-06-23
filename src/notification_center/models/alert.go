package models

type AlertVotes struct {
	IdVoter    string
	IdElection string
	MaxVotes   int
	Votes      int
	Emails     []string
}

type AlertCertificates struct {
	VoterId string
}
