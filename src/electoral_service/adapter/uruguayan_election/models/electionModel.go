package models

type ElectionJson struct {
	Error    bool          `json:"error"`
	Msg      string        `json:"msg"`
	Election ElectionModel `json:"election"`
}

type ElectionModel struct {
	Id            string   `bson:"id"`
	Description   string   `bson:"description"`
	StartingDate  string   `bson:"startingDate"`
	FinishingDate string   `bson:"finishingDate"`
	ElectionMode  string   `bson:"electionMode"`
	Emails        []string `json:"emails" bson:"emails"`

	Voters           []VoterModel          `json:"voters" bson:"voters"`
	Circuits         []CircuitModel        `json:"circuits"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
}
