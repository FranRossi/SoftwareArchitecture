package models

type ElectionModel struct {
	Id            string `faker:"customIdFaker"`
	Description   string `faker:"sentence"`
	StartingDate  string `faker:"date"`
	FinishingDate string `faker:"date"`
	ElectionMode  string `faker:"oneof: unico, multi"`

	Voter            []VoterModel          `faker:"custom_voter"`
	Candidate        []CandidateModel      `faker:"custom_candidates"`
	Circuit          []CircuitModel        `faker:"custom_circuits"`
	PoliticalParties []PoliticalPartyModel `faker:"custom_political_party"`
}
