package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/google/uuid"
)

type UserPreferencesRepository interface {
	Update(ctx context.Context, preferences *domain.UserPreferences) error

	List(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error)
}
