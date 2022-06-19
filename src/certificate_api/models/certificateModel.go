package models

type CertificateModel struct {
	IdVoter            string `bson:"id_voter"`
	IdElection         string `bson:"id_election"`
	TimeVoted          string `bson:"time_voted"`
	VoteIdentification string `bson:"vote_identification"`
	StartingDate       string `bson:"starting_date"`
	FinishingDate      string `bson:"finishing_date"`
	ElectionMode       string `bson:"election_mode"`
	Fullname           string `bson:"fullname"`
}
