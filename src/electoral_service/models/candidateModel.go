package models

type CandidateModel struct {
	Id                 string
	FullName           string
	Sex                string
	BirthDate          string
	IdPoliticalParty   string
	NamePoliticalParty string
	OtherFields        map[string]any `json:"otherFields" bson:"otherFields"`
}
