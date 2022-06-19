package jwt

import (
	"auth/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Manager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTManager(secretKey []byte, tokenDuration time.Duration) *Manager {
	return &Manager{secretKey, tokenDuration}
}

func (manager *Manager) Generate(user models.TokenInfo) (string, error) {
	claims := models.UserClaim{
		TokenInfo: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(manager.secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (manager *Manager) Verify(accessToken string) (*models.UserClaim, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return manager.secretKey, nil
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

type Roles struct {
	Consulter string
	Voter     string
	Electoral string
}

func GetRoles() Roles {
	return Roles{
		Consulter: "Consulter",
		Voter:     "Voter",
		Electoral: "Electoral",
	}
}
