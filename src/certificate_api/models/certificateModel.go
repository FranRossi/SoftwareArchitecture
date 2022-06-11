package models

type CertificateModel struct {
	VoterId       string `faker:"customIdFaker"`
	Description   string `faker:"sentence"`
	StartingDate  string `faker:"-"`
	FinishingDate string `faker:"-"`
	ElectionMode  string `faker:"oneof: unico, multi"`
	DateofVote 	  string
	Fullname	  string
	UniqueVoteId  string
}
