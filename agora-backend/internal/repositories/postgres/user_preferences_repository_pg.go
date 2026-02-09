package postgres

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPreferencesRepositoryPG struct {
	db *pgxpool.Pool
}

func NewUserPreferencesRepositoryPG(db *pgxpool.Pool) *UserPreferencesRepositoryPG {
	return &UserPreferencesRepositoryPG{db: db}
}

func (r *UserPreferencesRepositoryPG) Update(
	ctx context.Context,
	prefs *domain.UserPreferences,
) error {
	const cmd = `SELECT agora.update_user_preferences($1, $2);`

	tags := prefs.Tags()
	tagValues := make([]string, len(tags))
	for i, t := range tags {
		tagValues[i] = t.String()
	}

	_, err := r.db.Exec(ctx, cmd, prefs.UserID(), tagValues)
	return wrapPGError(err)
}

func (r *UserPreferencesRepositoryPG) List(
	ctx context.Context,
	userID uuid.UUID,
) (*domain.UserPreferences, error) {
	const cmd = `SELECT tag FROM agora.list_user_preferences($1);`

	rows, err := r.db.Query(ctx, cmd, userID)
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, wrapPGError(err)
		}

		tag, derr := domain.NewTag(value)
		if derr != nil {
			return nil, derr
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapPGError(err)
	}

	prefs, err := domain.NewUserPreferences(userID, tags)
	if err != nil {
		return nil, err
	}
	return prefs, nil
}
