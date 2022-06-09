package models

type Act struct {
	StarDate string `json:"startDate"`
	EndDate  string `json:"endDate"`
	Voters   int    `json:"voters"`
}
