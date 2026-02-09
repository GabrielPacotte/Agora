package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPG struct {
	db *pgxpool.Pool
}

func NewUserRepositoryPG(db *pgxpool.Pool) *UserRepositoryPG {
	return &UserRepositoryPG{db: db}
}

func (r *UserRepositoryPG) Create(ctx context.Context, user *domain.User, passwordHash string) error {
	cmd := `
        SELECT agora.create_user($1, $2, $3, $4, $5, $6, $7);
    `
	_, err := r.db.Exec(ctx, cmd,
		user.ID(),
		user.Login(),
		user.DisplayName(),
		passwordHash,
		true,
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	return wrapPGError(err)
}

func (r *UserRepositoryPG) Delete(ctx context.Context, userID uuid.UUID) error {
	cmd := `SELECT agora.delete_user($1);`
	_, err := r.db.Exec(ctx, cmd, userID)
	return wrapPGError(err)
}

func (r *UserRepositoryPG) UpdatePublicFields(ctx context.Context, user *domain.User) error {
	cmd := `
        SELECT agora.update_user_public_fields($1, $2, $3, $4, $5);
    `
	_, err := r.db.Exec(ctx, cmd,
		user.ID(),
		user.Login(),
		user.DisplayName(),
		true,
		user.UpdatedAt(),
	)
	return wrapPGError(err)
}

func (r *UserRepositoryPG) UpdatePassword(
	ctx context.Context,
	user *domain.User,
	newHash string,
) error {
	cmd := `SELECT agora.update_user_password($1, $2, $3);`
	_, err := r.db.Exec(ctx, cmd, user.ID(), newHash, user.UpdatedAt())
	return wrapPGError(err)
}

func (r *UserRepositoryPG) GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	cmd := `SELECT id, login, display_name, created_at, updated_at FROM agora.get_user_by_id($1);`

	row := r.db.QueryRow(ctx, cmd, userID)
	return scanUser(row)
}

func (r *UserRepositoryPG) GetByLoginWithPasswordHash(
	ctx context.Context,
	login string,
) (*domain.User, string, error) {
	const cmd = `
		SELECT id, login, display_name, password_hash, is_verified, created_at, updated_at
		FROM agora.get_user_by_login($1);
	`

	var (
		id           uuid.UUID
		dbLogin      string
		displayName  string
		passwordHash string
		isVerified   bool
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRow(ctx, cmd, login).Scan(
		&id,
		&dbLogin,
		&displayName,
		&passwordHash,
		&isVerified,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", wrapPGError(err)
	}

	user, derr := domain.ExistingUser(id, createdAt, updatedAt, dbLogin, displayName)
	if derr != nil {
		return nil, "", derr
	}

	return user, passwordHash, nil
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var (
		id        uuid.UUID
		login     string
		display   string
		createdAt time.Time
		updatedAt time.Time
	)

	if err := row.Scan(&id, &login, &display, &createdAt, &updatedAt); err != nil {
		return nil, wrapPGError(err)
	}

	return domain.ExistingUser(id, createdAt, updatedAt, login, display)
}
