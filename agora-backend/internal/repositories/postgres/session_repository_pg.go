package postgres

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepositoryPG struct {
	db *pgxpool.Pool
}

func NewSessionRepositoryPG(db *pgxpool.Pool) *SessionRepositoryPG {
	return &SessionRepositoryPG{db: db}
}

func (r *SessionRepositoryPG) Create(ctx context.Context, s *auth.Session, refreshTokenRaw string) error {
	cmd := `SELECT agora.create_session($1, $2, $3, $4, $5);`
	_, err := r.db.Exec(ctx, cmd,
		s.ID(),
		s.UserID(),
		refreshTokenRaw,
		s.CreatedAt(),
		s.ExpiresAt(),
	)
	return wrapPGError(err)
}

func (r *SessionRepositoryPG) GetByRefreshToken(ctx context.Context, refreshTokenRaw string) (*auth.Session, error) {
	cmd := `
		SELECT
			id,
			user_id,
			created_at,
			expires_at,
			revoked_at,
			last_used_at
		FROM agora.get_session_by_refresh_token($1);
	`
	row := r.db.QueryRow(ctx, cmd, refreshTokenRaw)
	return scanSession(row)
}

func (r *SessionRepositoryPG) TouchLastUsed(ctx context.Context, sessionID uuid.UUID, t time.Time) error {
	cmd := `SELECT agora.touch_session_last_used($1, $2);`
	_, err := r.db.Exec(ctx, cmd, sessionID, t)
	return wrapPGError(err)
}

func (r *SessionRepositoryPG) RevokeByRefreshToken(ctx context.Context, refreshTokenRaw string, t time.Time) error {
	cmd := `SELECT agora.revoke_session_by_refresh_token($1, $2);`
	_, err := r.db.Exec(ctx, cmd, refreshTokenRaw, t)
	return wrapPGError(err)
}

func scanSession(row pgx.Row) (*auth.Session, error) {
	var (
		id         uuid.UUID
		userID     uuid.UUID
		createdAt  time.Time
		expiresAt  time.Time
		revokedAt  *time.Time
		lastUsedAt *time.Time
	)

	if err := row.Scan(
		&id,
		&userID,
		&createdAt,
		&expiresAt,
		&revokedAt,
		&lastUsedAt,
	); err != nil {
		return nil, wrapPGError(err)
	}

	s, err := auth.ExistingSession(id, userID, createdAt, expiresAt, revokedAt, lastUsedAt)
	if err != nil {
		return nil, err
	}

	return s, nil
}
