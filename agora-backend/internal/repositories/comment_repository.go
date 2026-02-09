package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, commentID uuid.UUID) error

	GetByID(ctx context.Context, commentID uuid.UUID) (*domain.Comment, error)
	ListRootComments(ctx context.Context, postID uuid.UUID, page filters.Page) ([]domain.Comment, error)
	ListReplies(ctx context.Context, parentCommentID uuid.UUID, page filters.Page) ([]domain.Comment, error)
}
