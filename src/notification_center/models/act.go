package models

type InitialAct struct {
	StartDate        string                `json:"startDate"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	Voters           int                   `json:"voters"`
	Mode             string                `json:"mode"`
}

type ClosingAct struct {
	StarDate string         `json:"startDate"`
	EndDate  string         `json:"endDate"`
	Voters   int            `json:"voters"`
	Result   ResultElection `json:"result"`
}

type ResultElection struct {
	AmountVoted        int                        `json:"amount_voted"`
	VotesPerCandidates []CandidateEssential       `json:"votes_per_candidates"`
	VotesPerParties    []PoliticalPartyEssentials `json:"votes_per_parties"`
}

type PoliticalPartyEssentials struct {
	Name  string `bson:"name"`
	Votes int    `bson:"votes"`
}

type CandidateEssential struct {
	Id    string `bson:"id"`
	Name  string `bson:"name"`
	Votes int    `bson:"votes"`
}
