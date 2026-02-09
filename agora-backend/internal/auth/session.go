package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	id     uuid.UUID
	userID uuid.UUID

	createdAt  time.Time
	expiresAt  time.Time
	revokedAt  *time.Time
	lastUsedAt *time.Time
}

func ExistingSession(id, userID uuid.UUID, createdAt, expiresAt time.Time, revokedAt, lastUsedAt *time.Time) (*Session, error) {
	return &Session{
		id:         id,
		userID:     userID,
		createdAt:  createdAt,
		expiresAt:  expiresAt,
		revokedAt:  revokedAt,
		lastUsedAt: lastUsedAt,
	}, nil
}

func NewSession(id, userID uuid.UUID, createdAt, expiresAt time.Time) (*Session, error) {
	if id == uuid.Nil || userID == uuid.Nil || createdAt.IsZero() || expiresAt.IsZero() || expiresAt.Before(createdAt) {
		return nil, ErrInvalidSession
	}
	return &Session{id: id, userID: userID, createdAt: createdAt, expiresAt: expiresAt}, nil
}

func (s *Session) ID() uuid.UUID          { return s.id }
func (s *Session) UserID() uuid.UUID      { return s.userID }
func (s *Session) CreatedAt() time.Time   { return s.createdAt }
func (s *Session) ExpiresAt() time.Time   { return s.expiresAt }
func (s *Session) RevokedAt() *time.Time  { return s.revokedAt }
func (s *Session) LastUsedAt() *time.Time { return s.lastUsedAt }
func (s *Session) IsRevoked() bool        { return s.revokedAt != nil }
