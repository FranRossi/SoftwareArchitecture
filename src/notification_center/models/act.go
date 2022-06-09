package models

type Act struct {
	StarDate         string                `json:"startDate"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	EndDate          string                `json:"endDate"`
	Voters           int                   `json:"voters"`
	Mode             string                `json:"mode"`
}
