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
