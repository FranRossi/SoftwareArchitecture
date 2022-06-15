package models

import(
	jwt "github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	jwt.StandardClaims
	User `json:"user"`
}
