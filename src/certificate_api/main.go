package main

import (
	"certificate_api/connections"
)


func main {
	connections.Connection()

}

//funcion que recibe el modelo de info de la votacion y genera el certificado
type VoteInfo struct {
	IdVoter            string
	IdElection         string
	TimeVoted          string
	VoteIdentification string
}

func GenerateCertificate(voteinfo VoteInfo) error{
	var certificate models.CertificateModel
	certificate.IdVoter= voteinfo.IdVoter
	certificate.IdElection = voteinfo.IdElection 
	certificate.TimeVoted = voteinfo.TimeVoted 
	certificate.VoteIdentification= voteinfo.VoteIdentification
}