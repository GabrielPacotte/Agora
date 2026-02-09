package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type SavedContentRepository interface {
	Save(ctx context.Context, savedContent *domain.SavedContent) error
	Unsave(ctx context.Context, userID, subjectID uuid.UUID, contentType domain.ContentType) error

	List(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.SavedContent, error)
}
