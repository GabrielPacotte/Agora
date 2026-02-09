package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

// ---- NewReport ----

func TestNewReport_Valid(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	subjectID := uuid.New()

	desc := "   Something happened here   "

	r, err := domain.NewReport(now, authorID, subjectID, domain.POST, desc)
	if err != nil {
		t.Fatalf("NewReport unexpected error: %v", err)
	}

	if r == nil {
		t.Fatalf("NewReport returned nil")
	}

	// identified fields
	if r.ID() == uuid.Nil {
		t.Fatalf("expected non-nil report ID")
	}
	if !r.CreatedAt().Equal(now) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", now, r.CreatedAt())
	}
	if !r.UpdatedAt().Equal(now) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", now, r.UpdatedAt())
	}

	// IDs
	if r.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, r.AuthorID())
	}
	if r.SubjectID() != subjectID {
		t.Fatalf("SubjectID mismatch: expected %v, got %v", subjectID, r.SubjectID())
	}

	// subject type
	if r.SubjectType() != domain.POST {
		t.Fatalf("SubjectType mismatch: expected POST, got %v", r.SubjectType())
	}

	// description trimmed
	expectedDesc := strings.TrimSpace(desc)
	if r.Description() != expectedDesc {
		t.Fatalf("Description mismatch: expected %q, got %q", expectedDesc, r.Description())
	}

	if !r.HasDescription() {
		t.Fatalf("HasDescription should be true for a non-empty description")
	}
}

func TestNewReport_EmptyDescriptionAllowed(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	subjectID := uuid.New()

	r, err := domain.NewReport(now, authorID, subjectID, domain.COMMENT, "   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Description() != "" {
		t.Fatalf("expected empty trimmed description")
	}
	if r.HasDescription() {
		t.Fatalf("HasDescription should be false for empty description")
	}
}

func TestNewReport_InvalidIDs(t *testing.T) {
	now := time.Now()
	valid := uuid.New()

	tests := []struct {
		name      string
		authorID  uuid.UUID
		subjectID uuid.UUID
	}{
		{"nil authorID", uuid.Nil, valid},
		{"nil subjectID", valid, uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewReport(now, tt.authorID, tt.subjectID, domain.POST, "desc")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidID {
				t.Fatalf("expected ErrInvalidID, got %v", err)
			}
		})
	}
}

func TestNewReport_InvalidContentType(t *testing.T) {
	now := time.Now()

	_, err := domain.NewReport(now, uuid.New(), uuid.New(), "invalidType", "desc")
	if err == nil {
		t.Fatalf("expected error for invalid content type")
	}
	if err != domain.ErrInvalidContentType {
		t.Fatalf("expected ErrInvalidContentType, got %v", err)
	}
}

func TestNewReport_DescriptionTooLong(t *testing.T) {
	now := time.Now()
	longDesc := strings.Repeat("a", 2001) // > maxReportDescriptionLength

	_, err := domain.NewReport(now, uuid.New(), uuid.New(), domain.POST, longDesc)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrReportDescriptionTooLong {
		t.Fatalf("expected ErrReportDescriptionTooLong, got %v", err)
	}
}

func TestNewReport_InvalidTimestamp(t *testing.T) {
	_, err := domain.NewReport(time.Time{}, uuid.New(), uuid.New(), domain.POST, "desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidTimestamp {
		t.Fatalf("expected ErrInvalidTimestamp, got %v", err)
	}
}

// ---- ExistingReport ----

func TestExistingReport_Valid(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-1 * time.Hour)
	updated := time.Now()

	authorID := uuid.New()
	subjectID := uuid.New()
	desc := "  Something  "

	r, err := domain.ExistingReport(id, created, updated, authorID, subjectID, domain.COMMENT, desc)
	if err != nil {
		t.Fatalf("ExistingReport unexpected error: %v", err)
	}

	// ID + timestamps
	if r.ID() != id {
		t.Fatalf("ID mismatch: expected %v, got %v", id, r.ID())
	}
	if !r.CreatedAt().Equal(created) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", created, r.CreatedAt())
	}
	if !r.UpdatedAt().Equal(updated) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", updated, r.UpdatedAt())
	}

	// Data
	if r.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch")
	}
	if r.SubjectID() != subjectID {
		t.Fatalf("SubjectID mismatch")
	}
	if r.SubjectType() != domain.COMMENT {
		t.Fatalf("SubjectType mismatch")
	}
	if r.Description() != strings.TrimSpace(desc) {
		t.Fatalf("Description not trimmed correctly")
	}
}

func TestExistingReport_InvalidIDOrTimestamps(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()

	validAuthor := uuid.New()
	validSubject := uuid.New()

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
			_, err := domain.ExistingReport(tt.id, tt.created, tt.updated, validAuthor, validSubject, domain.POST, "desc")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestExistingReport_InvalidFields(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	created := now.Add(-1 * time.Hour)
	updated := now

	author := uuid.New()
	subject := uuid.New()

	// invalid content type
	_, err := domain.ExistingReport(id, created, updated, author, subject, "invalid", "desc")
	if err == nil {
		t.Fatalf("expected error for invalid content type")
	}
	if err != domain.ErrInvalidContentType {
		t.Fatalf("expected ErrInvalidContentType, got %v", err)
	}

	// description too long
	longDesc := strings.Repeat("x", 2001)

	_, err = domain.ExistingReport(id, created, updated, author, subject, domain.POST, longDesc)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrReportDescriptionTooLong {
		t.Fatalf("expected ErrReportDescriptionTooLong, got %v", err)
	}
}
