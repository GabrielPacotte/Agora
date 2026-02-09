package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type OpaqueTokenGenerator interface {
	GenerateOpaqueToken(nBytes int) (string, error)
}

type opaqueTokenGenerator struct{}

func (opaqueTokenGenerator) GenerateOpaqueToken(bytesLen int) (string, error) {
	if bytesLen < 32 {
		return "", ErrInvalidToken
	}
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func NewOpaqueTokenGenerator() opaqueTokenGenerator {
	return opaqueTokenGenerator{}
}
