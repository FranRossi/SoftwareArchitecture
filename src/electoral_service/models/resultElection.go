package models

type ResultElection struct {
	ElectionId          string                     `json:"election_id" bson:"election_id"`
	AmountOfVotes       int                        `json:"amount_voted" bson:"amount_voted"`
	TotalAmountOfVoters int                        `json:"voters" bson:"total_amount_of_voters"`
	VotesPerCandidates  []CandidateEssential       `json:"votes_per_candidates" bson:"votes_per_candidates"`
	VotesPerParties     []PoliticalPartyEssentials `json:"votes_per_parties" bson:"votes_per_parties"`
	Regions             []Region                   `json:"regions" bson:"regions"`
}

type Region struct {
	Name        string `json:"name" bson:"name"`
	TotalVoters int    `json:"total_voters" bson:"total_voters"`
	Votes       int    `json:"votes" bson:"votes"`
}
