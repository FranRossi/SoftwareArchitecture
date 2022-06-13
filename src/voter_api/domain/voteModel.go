package domain

type VoteModel struct {
	IdElection     string
	IdVoter        string
	Circuit        string
	IdCandidate    string
	PoliticalParty string
	Signature      []byte
}
