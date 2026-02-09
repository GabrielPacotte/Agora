package repositories

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error

	GetByID(ctx context.Context, reportID uuid.UUID) (*domain.Report, error)
	List(ctx context.Context, page filters.Page) ([]domain.Report, error)
	ListAgainstUser(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Report, error)
}
