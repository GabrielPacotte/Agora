package repositories

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/auth"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, s *auth.Session, refreshTokenRaw string) error
	GetByRefreshToken(ctx context.Context, refreshTokenRaw string) (*auth.Session, error)
	TouchLastUsed(ctx context.Context, sessionID uuid.UUID, t time.Time) error
	RevokeByRefreshToken(ctx context.Context, refreshTokenRaw string, t time.Time) error
}
