package jwt

import (
	"errors"
	"strconv"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwtv5.RegisteredClaims
}

type Manager struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

func NewManager(secret, issuer string, ttl time.Duration) *Manager {
	return &Manager{
		secretKey: []byte(secret),
		issuer:    issuer,
		ttl:       ttl,
	}
}

func (m *Manager) Generate(userID uint, email string) (string, *Claims, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(m.ttl)),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now),
			ID:        strconv.FormatInt(now.UnixNano(), 10),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", nil, err
	}
	return signed, claims, nil
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
