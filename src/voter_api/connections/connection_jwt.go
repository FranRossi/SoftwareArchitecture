package connections

import (
	jwt "auth/jwt"
	"time"
)

const (
	secretKey     = "secret"
	publicKey     = "public"
	tokenDuration = 15 * time.Minute
)

func ConnectionJWT() *jwt.Manager {
	return jwt.NewJWTManager([]byte(secretKey), []byte(publicKey), tokenDuration)
}
