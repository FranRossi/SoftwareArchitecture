package models

type ElectionResultJson struct {
	Error  bool          `json:"error"`
	Msg    string        `json:"msg"`
	Result ElectionModel `json:"request"`
}
type ElectionModel struct {
	Result              ResultElection `json:"result"`
	QueryRequestTime    string         `json:"query_request_time"`
	QueryResponseTime   string         `json:"query_response_time"`
	QueryProcessingTime string         `json:"query_processing_time"`
}

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

type PoliticalPartyEssentials struct {
	Name  string `bson:"name"`
	Votes int    `bson:"votes"`
}

type CandidateEssential struct {
	Id             string `bson:"id"`
	Name           string `bson:"name"`
	Votes          int    `bson:"votes"`
	PoliticalParty string `bson:"politicalParty"`
}

type ElectionConfigJson struct {
	Error  bool           `json:"error"`
	Msg    string         `json:"msg"`
	Config ElectionConfig `json:"request"`
}

type ElectionConfig struct {
	ElectionId          string   `json:"election_id"`
	MaxVotes            int      `json:"max_votes"`
	MaxCertificates     int      `json:"max_certificates"`
	Emails              []string `json:"emails"`
	QueryRequestTime    string   `json:"query_request_time"`
	QueryResponseTime   string   `json:"query_response_time"`
	QueryProcessingTime string   `json:"query_processing_time"`
}
