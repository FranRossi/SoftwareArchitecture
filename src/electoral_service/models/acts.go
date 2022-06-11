package models

type InitialAct struct {
	StarDate         string                `json:"startDate"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	Voters           int                   `json:"voters"`
	Mode             string                `json:"mode"`
}

type ClosingAct struct {
	StarDate            string         `json:"startDate"`
	EndDate             string         `json:"endDate"`
	TotalAmountOfVoters int            `json:"voters"`
	Result              ResultElection `json:"result"`
}

type ResultElection struct {
	AmountOfVotes      int                        `json:"amount_voted"`
	VotesPerCandidates []CandidateEssential       `json:"votes_per_candidates"`
	VotesPerParties    []PoliticalPartyEssentials `json:"votes_per_parties"`
}

type PoliticalPartyEssentials struct {
	Name  string `bson:"name"`
	Votes int    `bson:"votes"`
}

type CandidateEssential struct {
	Id             string `bson:"id"`
	Name           string `bson:"name"`
	Votes          int    `bson:"votes"`
	PoliticalParty string `bson:"politicalParty"`
}
