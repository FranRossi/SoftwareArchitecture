package models

type ElectionJson struct {
	Error    bool                   `json:"error"`
	Msg      string                 `json:"msg"`
	Election ElectionModelEssential `json:"election"`
}

type ElectionModelEssential struct {
	Id            string `bson:"id"`
	StartingDate  string `bson:"startingDate"`
	FinishingDate string `bson:"finishingDate"`
	ElectionMode  string `bson:"electionMode"` // "unico" o "multi"

	Voters           []VoterModel          `json:"voters" bson:"voters"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
	OtherFields      map[string]any        `json:"otherFields" bson:"otherFields"`
}
