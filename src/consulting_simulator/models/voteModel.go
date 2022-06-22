package models

type VoteModel struct {
	ElectionId          string `json:"electionId"`
	VoterId             string `json:"id"`
	TimeVoted           string `json:"time_voted"`
	QueryRequestTime    string `json:"query_request_time"`
	QueryResponseTime   string `json:"query_response_time"`
	QueryProcessingTime string `json:"query_processing_time"`
}

type VotesPerHours struct {
	ElectionId            string         `json:"electionId"`
	AmountOfVotesPerHours map[string]int `json:"amount_of_votes_per_hours"`
	QueryRequestTime      string         `json:"query_request_time"`
	QueryResponseTime     string         `json:"query_response_time"`
	QueryProcessingTime   string         `json:"query_processing_time"`
}

type VotesPerCircuits struct {
	ElectionId          string          `json:"electionId"`
	Circuit             string          `json:"circuit"`
	DataPerGroup        []VotesPerGroup `json:"data_per_group"`
	QueryRequestTime    string          `json:"query_request_time"`
	QueryResponseTime   string          `json:"query_response_time"`
	QueryProcessingTime string          `json:"query_processing_time"`
}

type VotesPerGroup struct {
	GroupName    string `json:"group_name"`
	MinAge       int    `json:"min_age"`
	MaxAge       int    `json:"max_age"`
	Sex          string `json:"sex"`
	CurrentVotes int    `json:"votes"`
	Total        int    `json:"capacity"`
}

type VotesPerRegion struct {
	ElectionId          string          `json:"electionId"`
	Region              string          `json:"region"`
	DataPerGroup        []VotesPerGroup `json:"data_per_group"`
	QueryRequestTime    string          `json:"query_request_time"`
	QueryResponseTime   string          `json:"query_response_time"`
	QueryProcessingTime string          `json:"query_processing_time"`
}

type VoteRegionCoverage struct {
	Error   bool             `json:"error"`
	Msg     string           `json:"msg"`
	Request VotesPerCircuits `json:"request"`
}

type VoteCircuitCoverage struct {
	Error   bool           `json:"error"`
	Msg     string         `json:"msg"`
	Request VotesPerRegion `json:"request"`
}
