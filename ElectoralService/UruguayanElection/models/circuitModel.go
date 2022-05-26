package models

type CircuitModel struct {
	Id         string `faker:"-"`
	Department string `faker:"-"`
	Address    string `faker:"sentence"`
}
