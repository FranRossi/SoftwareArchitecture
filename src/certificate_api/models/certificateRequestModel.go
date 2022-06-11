package models

type CertificateRequestModel struct {
	VoterId       string `faker:"customIdFaker"`
	Timestamp	  string
}