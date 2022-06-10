package models

type CandidateModel struct {
	Id               string
	Name             string
	LastName         string
	Sex              string
	BirthDate        string
	IdPoliticalParty string
	PoliticalParty   string
}

type CandidateEssential struct {
	Id             string `bson:"id"`
	Name           string `bson:"name"`
	Votes          int    `bson:"votes"`
	PoliticalParty string `bson:"politicalParty"`
}
