package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type ModerationActionRepository interface {
	Create(ctx context.Context, action *domain.ModerationAction) error
	Update(ctx context.Context, action *domain.ModerationAction) error

	List(ctx context.Context, page filters.Page) ([]domain.ModerationAction, error)
	ListAgainstUser(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.ModerationAction, error)
}
