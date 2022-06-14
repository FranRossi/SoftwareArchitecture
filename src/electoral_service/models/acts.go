package models

type InitialAct struct {
	StarDate         string                `json:"startDate"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	Voters           int                   `json:"voters"`
	Mode             string                `json:"mode"`
	ElectionId       string                `json:"electionId"`
}

type ClosingAct struct {
	StarDate   string         `json:"startDate"`
	EndDate    string         `json:"endDate"`
	Result     ResultElection `json:"result"`
	ElectionId string         `json:"electionId"`
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
