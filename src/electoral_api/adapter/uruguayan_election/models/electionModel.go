package models

type ElectionJson struct {
	Error    bool          `json:"error"`
	Msg      string        `json:"msg"`
	Election ElectionModel `json:"election"`
}

type ElectionModel struct {
	Id            string `bson:"id"`
	Description   string `bson:"description"`
	StartingDate  string `bson:"startingDate"`
	FinishingDate string `bson:"finishingDate"`
	ElectionMode  string `bson:"electionMode"`

	Voters           []VoterModel          `json:"voters" bson:"voters"`
	Candidates       []CandidateModel      `json:"candidates"`
	Circuits         []CircuitModel        `json:"circuits"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
}
