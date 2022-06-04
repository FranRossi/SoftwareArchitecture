package models

type ElectionJson struct {
	Error    bool          `json:"error"`
	Msg      string        `json:"msg"`
	Election ElectionModel `json:"election"`
}

type ElectionModel struct {
	Id            string `json:"id"`
	Description   string `json:"description"`
	StartingDate  string `json:"startingDate"`
	FinishingDate string `json:"finishingDate"`
	ElectionMode  string `json:"electionMode"`

	Voters           []VoterModel          `json:"voters"`
	Candidates       []CandidateModel      `json:"candidates"`
	Circuits         []CircuitModel        `json:"circuits"`
	PoliticalParties []PoliticalPartyModel `json:"politicalParties"`
}
