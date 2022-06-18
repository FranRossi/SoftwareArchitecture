package models

import(
	jwt "github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	jwt.StandardClaims
	TokenInfo `json:"user"`
}
