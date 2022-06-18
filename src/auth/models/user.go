package models


type TokenInfo struct {
	Id       string 
	Role     string 
}

type UserRegister struct {
	Id       string `json:"ci"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UserDB struct {
	Id string
	Role string
	HashedPassword string 

}