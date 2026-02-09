package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User, passwordHash string) error
	Delete(ctx context.Context, userID uuid.UUID) error
	UpdatePublicFields(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, user *domain.User, newHash string) error
	GetByLoginWithPasswordHash(ctx context.Context, login string) (*domain.User, string, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
