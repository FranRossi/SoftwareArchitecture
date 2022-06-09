package models

type PoliticalPartyModel struct {
	Id         string           `faker:"-"`
	Name       string           `faker:"-"`
	Candidates []CandidateModel `faker:"-"`
}
