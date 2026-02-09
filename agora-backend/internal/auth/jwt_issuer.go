package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTIssuer interface {
	GenerateAccessToken(userID uuid.UUID, now time.Time, ttl time.Duration) (AccessToken, error)
}

type jwtIssuer struct {
	secret   []byte
	issuer   string
	audience string
}

func NewJWTIssuer(secret string, issuer string) *jwtIssuer {
	return &jwtIssuer{
		secret: []byte(secret),
		issuer: issuer,
	}
}

func (j *jwtIssuer) GenerateAccessToken(userID uuid.UUID, now time.Time, ttl time.Duration) (AccessToken, error) {
	if userID == uuid.Nil {
		return "", errors.New("auth: nil userID")
	}
	if now.IsZero() {
		return "", errors.New("auth: zero now")
	}
	if ttl <= 0 {
		return "", errors.New("auth: ttl must be > 0")
	}

	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	raw, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return AccessToken(raw), nil
}
