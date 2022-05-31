package domain

type User struct {
	Id             string
	Username       string
	HashedPassword string
	Role           string
}
