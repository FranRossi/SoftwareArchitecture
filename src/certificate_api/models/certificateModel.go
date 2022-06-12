package models

type CertificateModel struct {
	IdVoter       string 
    IdElection         string
    TimeVoted          string
    VoteIdentification string

	StartingDate  string 
	FinishingDate string 
	ElectionMode  string 
	Fullname	  string
}
