package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, postID uuid.UUID) error

	ListByUserID(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error)
	GetFeed(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error)
	GetByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error)
	Search(ctx context.Context, filters filters.SearchPostFilter) ([]domain.Post, error)
}
