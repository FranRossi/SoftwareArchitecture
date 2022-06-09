package models

type PoliticalPartyModel struct {
	Id         string
	Name       string
	Candidates []CandidateModel
}
