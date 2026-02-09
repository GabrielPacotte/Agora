package services

import (
	"context"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type ModerationActionService interface {
	CreateAction(ctx context.Context, requesterID uuid.UUID, decision domain.ModerationDecision, reasonIDs []uuid.UUID, desc string) (*domain.ModerationAction, error)
	UpdateAction(ctx context.Context, requesterID, decisionID uuid.UUID, decision domain.ModerationDecision, reasonIDs []uuid.UUID, desc string) (*domain.ModerationAction, error)

	List(ctx context.Context, requesterID uuid.UUID, page filters.Page) ([]domain.ModerationAction, error)
	ListAgainstUser(ctx context.Context, requesterID, userID uuid.UUID, page filters.Page) ([]domain.ModerationAction, error)
}
