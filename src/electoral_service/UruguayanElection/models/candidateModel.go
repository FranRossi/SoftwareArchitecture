package models

type CandidateModel struct {
	Id               string `faker:"-"`
	Name             string `faker:"name"`
	LastName         string `faker:"last_name"`
	Sex              string `faker:"oneof: F,M"`
	BirthDate        string `faker:"date"`
	IdPoliticalParty string `faker:"-"`
}
