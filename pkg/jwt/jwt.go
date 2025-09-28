package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Manager struct {
	secret     string
	expiration time.Duration
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func NewManager(secret string, expiration time.Duration) *Manager {
	return &Manager{
		secret:     secret,
		expiration: expiration,
	}
}

func (m *Manager) GenerateToken(login string) (string, error) {
	claims := Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
