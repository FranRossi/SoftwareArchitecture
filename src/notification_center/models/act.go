package models

type InitialAct struct {
	StarDate         string                `json:"startDate"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	Voters           int                   `json:"voters"`
	Mode             string                `json:"mode"`
}

type ClosingAct struct {
	StarDate   string `json:"startDate"`
	EndDate    string `json:"endDate"`
	Voters     int    `json:"voters"`
	TotalVotes int    `json:"totalVotes"`
	Result     string `json:"result"`
}
