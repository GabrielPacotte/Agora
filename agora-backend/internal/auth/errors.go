package auth

import "errors"

var (
	ErrEmptyToken       = errors.New("empty token")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidExpiry    = errors.New("invalid expiry")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidTokenTTL  = errors.New("invalid token TTL")
	ErrInvalidSession   = errors.New("invalid session")
)
