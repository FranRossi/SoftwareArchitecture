package models

type PoliticalPartyModel struct {
	Id         string
	Name       string
	Candidates []CandidateModel
}

type CandidateModel struct {
	Id               string
	Name             string
	LastName         string
	Sex              string
	BirthDate        string
	IdPoliticalParty string
}
