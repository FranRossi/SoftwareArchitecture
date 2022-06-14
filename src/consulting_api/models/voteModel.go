package models

type VoteModel struct {
	VoterId             string `json:"id"`
	TimeVoted           string `json:"time_voted"`
	QueryRequestTime    string `json:"query_request_time"`
	QueryResponseTime   string `json:"query_response_time"`
	QueryProcessingTime string `json:"query_processing_time"`
}
