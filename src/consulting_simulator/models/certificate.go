package models

type CertificateJson struct {
	Error       bool                    `json:"error"`
	Msg         string                  `json:"msg"`
	Certificate CertificateRequestModel `json:"request"`
}

type CertificateRequestModel struct {
	VoterId            string `json:"voterId"`
	VoteIdentification string `json:"voteIdentification"`
	Timestamp          string `json:"timestamp"`
}

type CertificateResponseModel struct {
	IdVoter            string `bson:"id_voter"`
	IdElection         string `bson:"id_election"`
	TimeVoted          string `bson:"time_voted"`
	VoteIdentification string `bson:"vote_identification"`
	StartingDate       string `bson:"starting_date"`
	FinishingDate      string `bson:"finishing_date"`
	ElectionMode       string `bson:"election_mode"`
	Fullname           string `bson:"fullname"`
	Message            string `bson:"message"`
}
