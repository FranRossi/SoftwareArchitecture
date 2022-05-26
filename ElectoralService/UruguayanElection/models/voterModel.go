package models

type VoterModel struct {
	Id              string `faker:"-"`
	CivicCredential string `faker:"-"`
	Name            string `faker:"name"`
	LastName        string `faker:"last_name"`
	Sex             string `faker:"oneof: F, M"`
	BirthDate       string `faker:"date"`
	Department      string `faker:"-"`
	IdCircuit       string `faker:"-"`
	Phone           string `faker:"phone_number"`
	Email           string `faker:"email"`
}
