package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/repositories/postgres"
	"github.com/google/uuid"
)

func createTestUser(t *testing.T, ctx context.Context) *domain.User {
	t.Helper()

	userRepo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()
	login := "prefuser_" + uuid.New().String()[:8]

	user, err := domain.NewUser(now, login, "Prefs User")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	if err := userRepo.Create(ctx, user, "hash123"); err != nil {
		t.Fatalf("Create(user) error = %v", err)
	}

	return user
}

func TestUserPreferencesRepositoryPG_UpdateAndList(t *testing.T) {
	if testPool == nil {
		t.Skip("testPool is nil, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserPreferencesRepositoryPG(testPool)

	user := createTestUser(t, ctx)

	tagNews, err := domain.NewTag("news")
	if err != nil {
		t.Fatalf("NewTag(news) error = %v", err)
	}
	tagPolitics, err := domain.NewTag("politics")
	if err != nil {
		t.Fatalf("NewTag(politics) error = %v", err)
	}
	tagTech, err := domain.NewTag("tech")
	if err != nil {
		t.Fatalf("NewTag(tech) error = %v", err)
	}

	prefs, err := domain.NewUserPreferences(user.ID(), []domain.Tag{tagNews, tagPolitics, tagTech})
	if err != nil {
		t.Fatalf("NewUserPreferences() error = %v", err)
	}

	if err := repo.Update(ctx, prefs); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := repo.List(ctx, user.ID())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got == nil {
		t.Fatalf("List() returned nil preferences")
	}

	if got.UserID() != user.ID() {
		t.Errorf("UserID mismatch: got %v, want %v", got.UserID(), user.ID())
	}

	gotTags := got.Tags()
	if len(gotTags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(gotTags))
	}

	expected := map[string]struct{}{
		tagNews.String():     {},
		tagPolitics.String(): {},
		tagTech.String():     {},
	}

	for _, ttag := range gotTags {
		if _, ok := expected[ttag.String()]; !ok {
			t.Errorf("unexpected tag in preferences: %q", ttag.String())
		}
		delete(expected, ttag.String())
	}

	if len(expected) != 0 {
		t.Errorf("some expected tags are missing in results: %+v", expected)
	}
}

func TestUserPreferencesRepositoryPG_UpdateClearsWhenEmpty(t *testing.T) {
	if testPool == nil {
		t.Skip("testPool is nil, skipping integration test")
	}

	ctx := context.Background()
	repo := postgres.NewUserPreferencesRepositoryPG(testPool)

	user := createTestUser(t, ctx)

	tagNews, err := domain.NewTag("news")
	if err != nil {
		t.Fatalf("NewTag(news) error = %v", err)
	}
	initialPrefs, err := domain.NewUserPreferences(user.ID(), []domain.Tag{tagNews})
	if err != nil {
		t.Fatalf("NewUserPreferences(initial) error = %v", err)
	}

	if err := repo.Update(ctx, initialPrefs); err != nil {
		t.Fatalf("Update(initial) error = %v", err)
	}

	got, err := repo.List(ctx, user.ID())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.Tags()) != 1 {
		t.Fatalf("expected 1 tag after initial Update, got %d", len(got.Tags()))
	}

	emptyPrefs, err := domain.NewUserPreferences(user.ID(), nil)
	if err != nil {
		t.Fatalf("NewUserPreferences(empty) error = %v", err)
	}
	if err := repo.Update(ctx, emptyPrefs); err != nil {
		t.Fatalf("Update(empty) error = %v", err)
	}

	got2, err := repo.List(ctx, user.ID())
	if err != nil {
		t.Fatalf("List() after empty update error = %v", err)
	}
	if got2 == nil {
		t.Fatalf("List() returned nil preferences after empty update")
	}
	if len(got2.Tags()) != 0 {
		t.Fatalf("expected 0 tags after empty Update, got %d", len(got2.Tags()))
	}
}
