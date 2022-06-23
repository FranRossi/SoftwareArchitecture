package models

type UserRegister struct {
	Id       string `json:"ci"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type Login struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
