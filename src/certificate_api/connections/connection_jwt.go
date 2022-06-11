package connections

import (
	jwt "auth"
	"time"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func ConnectionJWT() *jwt.Manager {
	return jwt.NewJWTManager(secretKey, tokenDuration)
}
