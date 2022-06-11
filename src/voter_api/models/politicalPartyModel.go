package models

type PoliticalPartyModel struct {
	Id          string
	Name        string
	Candidates  []CandidateModel
	OtherFields map[string]any `json:"otherFields" bson:"otherFields"`
}
