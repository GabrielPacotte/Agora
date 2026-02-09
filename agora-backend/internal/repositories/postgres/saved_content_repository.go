package postgres

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SavedContentRepositoryPG struct {
	db *pgxpool.Pool
}

func NewSavedContentRepositoryPG(db *pgxpool.Pool) *SavedContentRepositoryPG {
	return &SavedContentRepositoryPG{db: db}
}

func (r *SavedContentRepositoryPG) Save(ctx context.Context, saved *domain.SavedContent) error {
	cmd := `
		SELECT agora.save_content($1, $2, $3);
	`
	_, err := r.db.Exec(ctx, cmd,
		saved.UserID(),
		saved.SubjectID(),
		saved.SubjectType(),
	)
	return wrapPGError(err)
}

func (r *SavedContentRepositoryPG) Unsave(
	ctx context.Context,
	userID, subjectID uuid.UUID,
	contentType domain.ContentType,
) error {
	cmd := `
		SELECT agora.unsave_content($1, $2, $3);
	`
	_, err := r.db.Exec(ctx, cmd, userID, subjectID, contentType)
	return wrapPGError(err)
}

func (r *SavedContentRepositoryPG) List(
	ctx context.Context,
	userID uuid.UUID,
	page filters.Page,
) ([]domain.SavedContent, error) {
	cmd := `
		SELECT user_id, content_id, content_type, created_at
		FROM agora.list_saved_content($1, $2, $3);
	`

	rows, err := r.db.Query(ctx, cmd, userID, page.Limit(), page.Offset())
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	var out []domain.SavedContent

	for rows.Next() {
		sc, err := scanSavedContent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *sc)
	}

	if rows.Err() != nil {
		return nil, wrapPGError(rows.Err())
	}

	return out, nil
}

func scanSavedContent(row pgx.Row) (*domain.SavedContent, error) {
	var (
		userID      uuid.UUID
		contentID   uuid.UUID
		contentType domain.ContentType
		createdAt   time.Time
	)

	if err := row.Scan(&userID, &contentID, &contentType, &createdAt); err != nil {
		return nil, wrapPGError(err)
	}

	return domain.NewSavedContent(createdAt, userID, contentID, contentType)
}
