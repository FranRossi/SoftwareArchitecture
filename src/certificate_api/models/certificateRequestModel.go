package models

type CertificateRequestModel struct {
	VoterId            string `json:"voterId"`
	VoteIdentification string `json:"voteIdentification"`
	Timestamp          string `json:"timestamp"`
}
