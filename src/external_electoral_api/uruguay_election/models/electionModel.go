package models

type ElectionModel struct {
	Id               string                `faker:"customIdFaker"`
	Description      string                `faker:"sentence"`
	StartingDate     string                `faker:"-"`
	FinishingDate    string                `faker:"-"`
	ElectionMode     string                `faker:"oneof: multi" //"oneof: unico, multi"`
	Emails           []string              `faker:"-"`
	Voters           []VoterModel          `faker:"custom_voter"`
	Circuits         []CircuitModel        `faker:"custom_circuits"`
	PoliticalParties []PoliticalPartyModel `faker:"custom_political_party"`
}
