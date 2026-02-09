package services

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type ReportService interface {
	ReportContent(ctx context.Context, authorID, contentID uuid.UUID, contentType domain.ContentType, description string) (*domain.Report, error)

	GetReport(ctx context.Context, requesterID, reportID uuid.UUID) (*domain.Report, error)
	ListReports(ctx context.Context, requesterID uuid.UUID, page filters.Page) ([]domain.Report, error)
	ListReportsAgainstUser(ctx context.Context, requesterID, userID uuid.UUID, page filters.Page) ([]domain.Report, error)
}
