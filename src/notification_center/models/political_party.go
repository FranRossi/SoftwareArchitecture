package models

type PoliticalPartyModel struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Candidates []CandidateModel `json:"candidates"`
}

type CandidateModel struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	LastName         string `json:"lastName"`
	Sex              string `json:"sex"`
	BirthDate        string `json:"birthDate"`
	IdPoliticalParty string `json:"idPoliticalParty"`
}
