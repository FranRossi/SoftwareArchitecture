package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"auth/models"
)

type Manager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *Manager {
	return &Manager{secretKey, tokenDuration}
}

func (manager *Manager) Generate(user models.TokenInfo) (string, error){
	claims := models.UserClaim{
		TokenInfo : user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}


func (manager *Manager) Verify(accessToken string) (*models.UserClaim, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*models.UserClaim)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}