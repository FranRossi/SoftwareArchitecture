package models

type VoterModel struct {
	Id              string `bson:"id"`
	CivicCredential string `bson:"civicCredential"`
	Name            string `bson:"name"`
	LastName        string `bson:"lastName"`
	Sex             string `bson:"sex"`
	BirthDate       string `bson:"birthDate"`
	Department      string `bson:"department"`
	IdCircuit       string `bson:"idCircuit"`
	Phone           string `bson:"phone"`
	Email           string `bson:"email"`
	Voted           int    `bson:"voted"`
}
