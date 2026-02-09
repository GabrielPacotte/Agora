package auth

import "time"

type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
)

type TokenPair struct {
	Access    AccessToken
	Refresh   RefreshToken
	Type      TokenType
	ExpiresIn time.Duration
}

func NewTokenPair(
	access AccessToken,
	refresh RefreshToken,
	accessTTL time.Duration,
) (TokenPair, error) {
	if access.IsZero() {
		return TokenPair{}, ErrInvalidToken
	}
	if refresh.IsZero() || refresh.ExpiresAt().IsZero() {
		return TokenPair{}, ErrInvalidToken
	}
	if accessTTL <= 0 {
		return TokenPair{}, ErrInvalidTokenTTL
	}

	return TokenPair{
		Type:      TokenTypeBearer,
		Access:    access,
		Refresh:   refresh,
		ExpiresIn: accessTTL,
	}, nil
}
