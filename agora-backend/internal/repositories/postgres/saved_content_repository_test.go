package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/repositories/postgres"
	"github.com/google/uuid"
)

func TestSavedContentRepositoryPG_SaveAndList(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	userRepo := postgres.NewUserRepositoryPG(testPool)
	postRepo := postgres.NewPostRepositoryPG(testPool)
	savedRepo := postgres.NewSavedContentRepositoryPG(testPool)

	now := time.Now().UTC()

	// --- user ---
	user, err := domain.NewUser(now, "saved_user", "Saved User")
	if err != nil {
		t.Fatalf("NewUser error: %v", err)
	}
	if err := userRepo.Create(ctx, user, "hash123"); err != nil {
		t.Fatalf("userRepo.Create error: %v", err)
	}

	// --- post + stances ---
	color1, _ := domain.NewColor("#ff0000")
	color2, _ := domain.NewColor("#00ff00")
	color3, _ := domain.NewColor("#0000ff")

	s1, _ := domain.NewStance(uuid.New(), "Agree", color1, "")
	s2, _ := domain.NewStance(uuid.New(), "Disagree", color2, "")
	s3, _ := domain.NewStance(uuid.New(), "Neutral", color3, "")

	post, err := domain.NewPost(
		now,
		"Post to save",
		"Body",
		nil,
		[]domain.Stance{s1, s2, s3},
		user.ID(),
	)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("postRepo.Create error: %v", err)
	}

	saved, err := domain.NewSavedContent(
		now,
		user.ID(),
		post.ID(),
		domain.POST,
	)
	if err != nil {
		t.Fatalf("NewSavedContent error: %v", err)
	}

	if err := savedRepo.Save(ctx, saved); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	page := filters.NewPage(10, 0)
	list, err := savedRepo.List(ctx, user.ID(), page)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 saved content, got %d", len(list))
	}

	got := list[0]

	if got.UserID() != user.ID() {
		t.Errorf("UserID mismatch: got %v, want %v", got.UserID(), user.ID())
	}
	if got.SubjectID() != post.ID() {
		t.Errorf("ContentID mismatch: got %v, want %v", got.SubjectID(), post.ID())
	}
	if got.SubjectType() != domain.POST {
		t.Errorf("ContentType mismatch: got %v, want %v", got.SubjectID(), domain.POST)
	}
	if got.SavedAt().IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
}

func TestSavedContentRepositoryPG_Unsave(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	userRepo := postgres.NewUserRepositoryPG(testPool)
	postRepo := postgres.NewPostRepositoryPG(testPool)
	savedRepo := postgres.NewSavedContentRepositoryPG(testPool)

	now := time.Now().UTC()

	// user
	user, err := domain.NewUser(now, "unsave_user", "Unsave User")
	if err != nil {
		t.Fatalf("NewUser error: %v", err)
	}
	if err := userRepo.Create(ctx, user, "hash123"); err != nil {
		t.Fatalf("userRepo.Create error: %v", err)
	}

	// post
	color1, _ := domain.NewColor("#ff0000")
	color2, _ := domain.NewColor("#00ff00")
	color3, _ := domain.NewColor("#0000ff")

	s1, _ := domain.NewStance(uuid.New(), "Agree", color1, "")
	s2, _ := domain.NewStance(uuid.New(), "Disagree", color2, "")
	s3, _ := domain.NewStance(uuid.New(), "Neutral", color3, "")

	post, err := domain.NewPost(
		now,
		"Post to unsave",
		"Body",
		nil,
		[]domain.Stance{s1, s2, s3},
		user.ID(),
	)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("postRepo.Create error: %v", err)
	}

	// save
	saved, err := domain.NewSavedContent(now, user.ID(), post.ID(), domain.POST)
	if err != nil {
		t.Fatalf("NewSavedContent error: %v", err)
	}
	if err := savedRepo.Save(ctx, saved); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// unsave
	if err := savedRepo.Unsave(ctx, user.ID(), post.ID(), domain.POST); err != nil {
		t.Fatalf("Unsave error: %v", err)
	}

	// check list is empty
	page := filters.NewPage(10, 0)
	list, err := savedRepo.List(ctx, user.ID(), page)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 saved content after unsave, got %d", len(list))
	}
}

func TestSavedContentRepositoryPG_Unsave_NotFound(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	savedRepo := postgres.NewSavedContentRepositoryPG(testPool)

	userID := uuid.New()
	contentID := uuid.New()

	err := savedRepo.Unsave(ctx, userID, contentID, domain.POST)
	if err == nil {
		t.Fatalf("expected error on unsave non-existing saved content, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
