package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

// ---- NewModerationAction ----

func TestNewModerationAction_Valid(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	r1 := uuid.New()
	r2 := uuid.New()

	reasons := []uuid.UUID{r1, r2, r1}
	desc := "   Action taken on user   "

	a, err := domain.NewModerationAction(now, authorID, domain.ContentDeleted, reasons, desc)
	if err != nil {
		t.Fatalf("NewModerationAction unexpected error: %v", err)
	}
	if a == nil {
		t.Fatalf("NewModerationAction returned nil")
	}

	if a.ID() == uuid.Nil {
		t.Fatalf("expected non-nil ID")
	}
	if !a.CreatedAt().Equal(now) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", now, a.CreatedAt())
	}
	if !a.UpdatedAt().Equal(now) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", now, a.UpdatedAt())
	}

	if a.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, a.AuthorID())
	}

	if a.Decision() != domain.ContentDeleted {
		t.Fatalf("Decision mismatch: expected ContentDeleted, got %v", a.Decision())
	}

	expectedDesc := strings.TrimSpace(desc)
	if a.Description() != expectedDesc {
		t.Fatalf("Description mismatch: expected %q, got %q", expectedDesc, a.Description())
	}

	gotReasons := a.ReasonIDs()
	if len(gotReasons) != 2 {
		t.Fatalf("expected 2 unique reasons, got %d", len(gotReasons))
	}
	if gotReasons[0] != r1 || gotReasons[1] != r2 {
		t.Fatalf("unexpected reasons order/content: %v", gotReasons)
	}
}

func TestNewModerationAction_InvalidAuthorID(t *testing.T) {
	now := time.Now()
	_, err := domain.NewModerationAction(now, uuid.Nil, domain.ReportReviewed, []uuid.UUID{uuid.New()}, "reason")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestNewModerationAction_InvalidDecision(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	reasons := []uuid.UUID{uuid.New()}

	invalidDecision := domain.ModerationDecision("not-valid")

	_, err := domain.NewModerationAction(now, authorID, invalidDecision, reasons, "desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidModerationActionDecision {
		t.Fatalf("expected ErrInvalidModerationActionDecision, got %v", err)
	}
}

func TestNewModerationAction_EmptyReasons(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	_, err := domain.NewModerationAction(now, authorID, domain.ReportReviewed, nil, "desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrEmptyReasons {
		t.Fatalf("expected ErrEmptyReasons, got %v", err)
	}
}

func TestNewModerationAction_InvalidReasonIDs(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	reasons := []uuid.UUID{uuid.New(), uuid.Nil}

	_, err := domain.NewModerationAction(now, authorID, domain.ReportReviewed, reasons, "desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestNewModerationAction_InvalidDescription(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	reasons := []uuid.UUID{uuid.New()}

	_, err := domain.NewModerationAction(now, authorID, domain.ReportReviewed, reasons, "   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrModerationActionDescriptionSize {
		t.Fatalf("expected ErrModerationActionDescriptionSize, got %v", err)
	}

	longDesc := strings.Repeat("x", 1001)
	_, err = domain.NewModerationAction(now, authorID, domain.ReportReviewed, reasons, longDesc)
	if err == nil {
		t.Fatalf("expected error for long description, got nil")
	}
	if err != domain.ErrModerationActionDescriptionSize {
		t.Fatalf("expected ErrModerationActionDescriptionSize, got %v", err)
	}
}

func TestNewModerationAction_InvalidTimestamp(t *testing.T) {
	_, err := domain.NewModerationAction(time.Time{}, uuid.New(), domain.ReportReviewed, []uuid.UUID{uuid.New()}, "desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidTimestamp {
		t.Fatalf("expected ErrInvalidTimestamp, got %v", err)
	}
}

// ---- ExistingModerationAction ----

func TestExistingModerationAction_Valid(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	authorID := uuid.New()
	r1 := uuid.New()
	r2 := uuid.New()

	reasons := []uuid.UUID{r1, r2, r1}
	desc := "   Some decision   "

	a, err := domain.ExistingModerationAction(id, created, updated, authorID, domain.UserSuspended, reasons, desc)
	if err != nil {
		t.Fatalf("ExistingModerationAction unexpected error: %v", err)
	}

	if a.ID() != id {
		t.Fatalf("ID mismatch: expected %v, got %v", id, a.ID())
	}
	if !a.CreatedAt().Equal(created) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", created, a.CreatedAt())
	}
	if !a.UpdatedAt().Equal(updated) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", updated, a.UpdatedAt())
	}

	if a.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, a.AuthorID())
	}
	if a.Decision() != domain.UserSuspended {
		t.Fatalf("Decision mismatch: expected UserSuspended, got %v", a.Decision())
	}

	if a.Description() != strings.TrimSpace(desc) {
		t.Fatalf("Description not trimmed correctly")
	}

	gotReasons := a.ReasonIDs()
	if len(gotReasons) != 2 {
		t.Fatalf("expected 2 unique reasons, got %d", len(gotReasons))
	}
	if gotReasons[0] != r1 || gotReasons[1] != r2 {
		t.Fatalf("unexpected reasons order/content: %v", gotReasons)
	}
}

func TestExistingModerationAction_InvalidAuthorOrDecisionOrReasons(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	validAuthor := uuid.New()
	validReasons := []uuid.UUID{uuid.New()}
	longDesc := strings.Repeat("x", 1001)
	nilReason := []uuid.UUID{uuid.Nil}
	invalidDecision := domain.ModerationDecision("not-valid")

	tests := []struct {
		name     string
		authorID uuid.UUID
		decision domain.ModerationDecision
		reasons  []uuid.UUID
		desc     string
		wantErr  error
	}{
		{"nil authorID", uuid.Nil, domain.ReportReviewed, validReasons, "desc", domain.ErrInvalidID},
		{"invalid decision", validAuthor, invalidDecision, validReasons, "desc", domain.ErrInvalidModerationActionDecision},
		{"no reasons", validAuthor, domain.ReportReviewed, nil, "desc", domain.ErrEmptyReasons},
		{"nil reason ID", validAuthor, domain.ReportReviewed, nilReason, "desc", domain.ErrInvalidID},
		{"empty description after trim", validAuthor, domain.ReportReviewed, validReasons, "   ", domain.ErrModerationActionDescriptionSize},
		{"too long description", validAuthor, domain.ReportReviewed, validReasons, longDesc, domain.ErrModerationActionDescriptionSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ExistingModerationAction(id, created, updated, tt.authorID, tt.decision, tt.reasons, tt.desc)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestExistingModerationAction_InvalidIdentified(t *testing.T) {
	id := uuid.New()
	authorID := uuid.New()
	reasons := []uuid.UUID{uuid.New()}

	created := time.Now().Add(-time.Hour)
	updated := time.Now()

	tests := []struct {
		name    string
		id      uuid.UUID
		created time.Time
		updated time.Time
		wantErr error
	}{
		{"nil ID", uuid.Nil, created, updated, domain.ErrInvalidID},
		{"createdAt zero", id, time.Time{}, updated, domain.ErrInvalidTimestamp},
		{"updatedAt zero", id, created, time.Time{}, domain.ErrInvalidTimestamp},
		{"updated before created", id, updated, created, domain.ErrInvalidTimestamp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ExistingModerationAction(tt.id, tt.created, tt.updated, authorID, domain.ModerationPending, reasons, "desc")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestModerationAction_ReasonIDsReturnsCopy(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	r1 := uuid.New()
	r2 := uuid.New()

	a, err := domain.NewModerationAction(now, authorID, domain.ReportReviewed, []uuid.UUID{r1, r2}, "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reasons1 := a.ReasonIDs()
	reasons1[0] = uuid.New()

	reasons2 := a.ReasonIDs()
	if reasons2[0] != r1 {
		t.Fatalf("expected original first reason to remain %v, got %v", r1, reasons2[0])
	}
}
