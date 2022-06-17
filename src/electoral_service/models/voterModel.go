package models

type VoterModel struct {
	Id                   string         `bson:"id"`
	FullName             string         `bson:"name"`
	Sex                  string         `bson:"sex"`
	BirthDate            string         `bson:"birthDate"`
	Phone                string         `bson:"phone"`
	Email                string         `bson:"email"`
	Voted                int            `bson:"voted"`
	LastCandidateVotedId string         `bson:"lastCandidateVotedId"`
	Region               string         `bson:"region"`
	OtherFields          map[string]any `json:"otherFields" bson:"otherFields"`
}
