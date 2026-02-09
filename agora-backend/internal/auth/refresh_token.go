package auth

import (
	"strings"
	"time"
)

type RefreshToken struct {
	raw       string
	expiresAt time.Time
}

func NewRefreshToken(raw string, expiresAt time.Time) (RefreshToken, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return RefreshToken{}, ErrEmptyToken
	}
	if expiresAt.IsZero() {
		return RefreshToken{}, ErrInvalidExpiry
	}
	return RefreshToken{raw: raw, expiresAt: expiresAt}, nil
}

func (t RefreshToken) Raw() string          { return t.raw }
func (t RefreshToken) ExpiresAt() time.Time { return t.expiresAt }
func (t RefreshToken) IsZero() bool         { return t.raw == "" }
