package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

func TestNewUser_Valid(t *testing.T) {
	now := time.Now()

	login := "  John_Doe-42  "
	displayName := "  John Doe  "

	u, err := domain.NewUser(now, login, displayName)
	if err != nil {
		t.Fatalf("NewUser unexpected error: %v", err)
	}
	if u == nil {
		t.Fatalf("NewUser returned nil user")
	}

	// ID / timestamps (via identified embedding)
	if u.ID() == uuid.Nil {
		t.Fatalf("expected non-nil user ID")
	}
	if !u.CreatedAt().Equal(now) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", now, u.CreatedAt())
	}
	if !u.UpdatedAt().Equal(now) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", now, u.UpdatedAt())
	}

	// login: trim + lowercase
	expectedLogin := strings.ToLower(strings.TrimSpace(login))
	if u.Login() != expectedLogin {
		t.Fatalf("Login mismatch: expected %q, got %q", expectedLogin, u.Login())
	}

	// displayName: trim
	expectedDisplayName := strings.TrimSpace(displayName)
	if u.DisplayName() != expectedDisplayName {
		t.Fatalf("DisplayName mismatch: expected %q, got %q", expectedDisplayName, u.DisplayName())
	}
}

func TestNewUser_InvalidLogin(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		login string
	}{
		{"empty", ""},
		{"spaces only", "   "},
		{"too short", "a"},
		{"invalid chars", "john#doe"},
		{"starts with dash", "-johndoe"},
		{"ends with underscore", "johndoe_"},
		{"too long", strings.Repeat("a", 40)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewUser(now, tt.login, "John Doe")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidLogin {
				t.Fatalf("expected ErrInvalidLogin, got %v", err)
			}
		})
	}
}

func TestNewUser_InvalidDisplayName(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		displayName string
	}{
		{"empty", ""},
		{"spaces only", "    "},
		{"too long", strings.Repeat("x", 51)},
		{"control char", "John \n Doe"},
		{"invalid symbol", "John_Doe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewUser(now, "johndoe", tt.displayName)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidDisplayName {
				t.Fatalf("expected ErrInvalidDisplayName, got %v", err)
			}
		})
	}
}

func TestNewUser_InvalidTimestamp(t *testing.T) {
	_, err := domain.NewUser(time.Time{}, "johndoe", "John Doe")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidTimestamp {
		t.Fatalf("expected ErrInvalidTimestamp, got %v", err)
	}
}

// ---- ExistingUser ----

func TestExistingUser_Valid(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()

	login := "  JOHNDOE  "
	displayName := "  John Doe  "

	u, err := domain.ExistingUser(id, created, updated, login, displayName)
	if err != nil {
		t.Fatalf("ExistingUser unexpected error: %v", err)
	}

	if u.ID() != id {
		t.Fatalf("ID mismatch: expected %v, got %v", id, u.ID())
	}
	if !u.CreatedAt().Equal(created) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", created, u.CreatedAt())
	}
	if !u.UpdatedAt().Equal(updated) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", updated, u.UpdatedAt())
	}

	expectedLogin := strings.ToLower(strings.TrimSpace(login))
	if u.Login() != expectedLogin {
		t.Fatalf("Login mismatch: expected %q, got %q", expectedLogin, u.Login())
	}

	expectedDisplayName := strings.TrimSpace(displayName)
	if u.DisplayName() != expectedDisplayName {
		t.Fatalf("DisplayName mismatch: expected %q, got %q", expectedDisplayName, u.DisplayName())
	}
}

func TestExistingUser_InvalidIDOrTimestamps(t *testing.T) {
	id := uuid.New()
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
			_, err := domain.ExistingUser(tt.id, tt.created, tt.updated, "johndoe", "John Doe")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestExistingUser_InvalidLoginOrDisplayName(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()

	_, err := domain.ExistingUser(id, created, updated, "  ", "John Doe")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidLogin {
		t.Fatalf("expected ErrInvalidLogin, got %v", err)
	}

	_, err = domain.ExistingUser(id, created, updated, "johndoe", "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrInvalidDisplayName {
		t.Fatalf("expected ErrInvalidDisplayName, got %v", err)
	}
}
