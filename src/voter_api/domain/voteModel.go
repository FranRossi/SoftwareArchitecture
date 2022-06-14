package domain

type VoteModel struct {
	IdElection     string
	IdVoter        string
	Circuit        string
	IdCandidate    string
	PoliticalParty string
	Signature      []byte
}

type VoteInfo struct {
	IdVoter            string
	IdElection         string
	TimeVoted          string
	VoteIdentification string
}
