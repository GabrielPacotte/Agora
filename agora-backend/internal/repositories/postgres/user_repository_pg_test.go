package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/repositories/postgres"
	"github.com/GabrielPacotte/Agora/internal/testutils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func applySQLFile(ctx context.Context, pool *pgxpool.Pool, relPath string) error {
	root := testutils.ProjectRoot()
	fullPath := filepath.Join(root, relPath)

	sqlBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", fullPath, err)
	}

	_, err = pool.Exec(ctx, string(sqlBytes))
	if err != nil {
		return fmt.Errorf("executing %s: %w", fullPath, err)
	}

	return nil
}

func TestUserRepositoryPG_CreateAndGetByID(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	user, err := domain.NewUser(now, "johndoe", "John Doe")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	// Create
	if err := repo.Create(ctx, user, "hash123"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// GetByID
	got, err := repo.GetByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got == nil {
		t.Fatalf("GetByID() returned nil user")
	}

	if got.ID() != user.ID() {
		t.Errorf("ID mismatch: got %v, want %v", got.ID(), user.ID())
	}
	if got.Login() != user.Login() {
		t.Errorf("Login mismatch: got %q, want %q", got.Login(), user.Login())
	}
	if got.DisplayName() != user.DisplayName() {
		t.Errorf("DisplayName mismatch: got %q, want %q", got.DisplayName(), user.DisplayName())
	}

	if got.CreatedAt().IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
	if got.UpdatedAt().IsZero() {
		t.Errorf("UpdatedAt should not be zero")
	}
}

func TestUserRepositoryPG_Create_DuplicateLogin(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	user1, err := domain.NewUser(now, "dupuser", "Dup User 1")
	if err != nil {
		t.Fatalf("NewUser(user1) error = %v", err)
	}
	if err := repo.Create(ctx, user1, "hashpass"); err != nil {
		t.Fatalf("Create(user1) error = %v", err)
	}

	user2, err := domain.NewUser(now, "dupuser", "Dup User 2")
	if err != nil {
		t.Fatalf("NewUser(user2) error = %v", err)
	}

	err = repo.Create(ctx, user2, "hashpass")
	if err == nil {
		t.Fatalf("expected error on duplicate login, got nil")
	}

	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestUserRepositoryPG_GetByID_NotFound(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	randomID := uuid.New()

	user, err := repo.GetByID(ctx, randomID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user on not found, got %+v", user)
	}
}

func TestUserRepositoryPG_UpdatePublicFields_Success(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	user, err := domain.NewUser(now, "updateuser", "Old Name")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	if err := repo.Create(ctx, user, "hash-xyz"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := domain.ExistingUser(
		user.ID(),
		user.CreatedAt(),
		user.UpdatedAt(),
		"newlogin",
		"New Name",
	)
	if err != nil {
		t.Fatalf("ExistingUser() error = %v", err)
	}

	if err := repo.UpdatePublicFields(ctx, updated); err != nil {
		t.Fatalf("UpdatePublicFields() error = %v", err)
	}

	got, err := repo.GetByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.Login() != "newlogin" {
		t.Errorf("expected login %q, got %q", "newlogin", got.Login())
	}
	if got.DisplayName() != "New Name" {
		t.Errorf("expected displayName %q, got %q", "New Name", got.DisplayName())
	}
}

func TestUserRepositoryPG_UpdatePublicFields_DuplicateLogin(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	userA, err := domain.NewUser(now, "user_a", "User A")
	if err != nil {
		t.Fatalf("NewUser(A) error = %v", err)
	}
	if err := repo.Create(ctx, userA, "hash-a"); err != nil {
		t.Fatalf("Create(A) error = %v", err)
	}

	userB, err := domain.NewUser(now, "user_b", "User B")
	if err != nil {
		t.Fatalf("NewUser(B) error = %v", err)
	}
	if err := repo.Create(ctx, userB, "hash-b"); err != nil {
		t.Fatalf("Create(B) error = %v", err)
	}

	updatedB, err := domain.ExistingUser(
		userB.ID(),
		userB.CreatedAt(),
		userB.UpdatedAt(),
		"user_a",
		userB.DisplayName(),
	)
	if err != nil {
		t.Fatalf("ExistingUser(B updated) error = %v", err)
	}

	err = repo.UpdatePublicFields(ctx, updatedB)
	if err == nil {
		t.Fatalf("expected error on duplicate login, got nil")
	}
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestUserRepositoryPG_Delete_Success(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	user, err := domain.NewUser(now, "todelete", "To Delete")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	if err := repo.Create(ctx, user, "hash-del"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := repo.Delete(ctx, user.ID()); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, user.ID())
	if err == nil {
		t.Fatalf("expected not found after delete, got nil error")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestUserRepositoryPG_Delete_NotFound(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	randomID := uuid.New()

	err := repo.Delete(ctx, randomID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepositoryPG_UpdatePassword_Success(t *testing.T) {
	if testPool == nil {
		t.Skip("KAINE_TEST_DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	user, err := domain.NewUser(now, "pwduser", "Pwd User")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	if err := repo.Create(ctx, user, "old-hash"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = repo.UpdatePassword(ctx, user, "old-hash")
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}
}
