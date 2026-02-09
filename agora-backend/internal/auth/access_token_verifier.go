package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenVerifier interface {
	Verify(token string) (uuid.UUID, error)
}

type accessTokenVerifier struct {
	secret []byte
	issuer string
}

func NewAccessTokenVerifier(secret []byte, issuer string) *accessTokenVerifier {
	return &accessTokenVerifier{secret: secret, issuer: issuer}
}

func (v *accessTokenVerifier) Verify(raw string) (uuid.UUID, error) {
	token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return v.secret, nil
	}, jwt.WithIssuer(v.issuer))
	if err != nil || token == nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return id, nil
}
