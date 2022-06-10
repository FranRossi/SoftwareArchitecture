package models

type InitialAct struct {
	StarDate         string                `json:"startDate"`
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
	AmountVoted int                  `json:"amount_voted"`
	Candidates  []CandidateEssential `json:"candidates"`
}
