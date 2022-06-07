package models

type ElectionModel struct {
	Id            string `faker:"customIdFaker"`
	Description   string `faker:"sentence"`
	StartingDate  string `faker:"-"`
	FinishingDate string `faker:"-"`
	ElectionMode  string `faker:"oneof: unico, multi"`

	Voters           []VoterModel          `faker:"custom_voter"`
	Candidates       []CandidateModel      `faker:"custom_candidates"`
	Circuits         []CircuitModel        `faker:"custom_circuits"`
	PoliticalParties []PoliticalPartyModel `faker:"custom_political_party"`
}
