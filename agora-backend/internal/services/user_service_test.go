package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/services"
	"github.com/google/uuid"
)

func TestUserService_GetProfile(t *testing.T) {
	uid := uuid.New()

	mockUsers := &mockUserRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			if id == uid {
				u, _ := domain.NewUser(time.Now(), "john", "John Doe")
				return u, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewUserService(mockUsers, nil, nil)

	t.Run("success", func(t *testing.T) {
		u, err := svc.GetProfile(context.Background(), uid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Login() != "john" {
			t.Fatalf("expected login john, got %s", u.Login())
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetProfile(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestUserService_ListPreferences(t *testing.T) {
	uid := uuid.New()

	tag1, _ := domain.NewTag("politics")
	tag2, _ := domain.NewTag("science")

	prefs, _ := domain.NewUserPreferences(uid, []domain.Tag{tag1, tag2})

	mockPrefs := &mockUserPreferencesRepository{
		ListFn: func(ctx context.Context, id uuid.UUID) (*domain.UserPreferences, error) {
			if id == uid {
				return prefs, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewUserService(nil, mockPrefs, nil)

	t.Run("success", func(t *testing.T) {
		out, err := svc.ListPreferences(context.Background(), uid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) != 2 {
			t.Fatalf("expected 2 tags, got %d", len(out))
		}
	})

	t.Run("no preferences -> empty list", func(t *testing.T) {
		out, err := svc.ListPreferences(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) != 0 {
			t.Fatalf("expected empty slice, got %d items", len(out))
		}
	})
}

func TestUserService_SetPreferences(t *testing.T) {
	uid := uuid.New()
	tag1, _ := domain.NewTag("go")
	tag2, _ := domain.NewTag("dev")

	captured := []*domain.UserPreferences{}

	mockPrefs := &mockUserPreferencesRepository{
		UpdateFn: func(ctx context.Context, prefs *domain.UserPreferences) error {
			captured = append(captured, prefs)
			return nil
		},
	}

	svc := services.NewUserService(nil, mockPrefs, nil)

	err := svc.SetPreferences(context.Background(), uid, []domain.Tag{tag1, tag2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(captured) != 1 {
		t.Fatalf("expected 1 update call, got %d", len(captured))
	}

	gotPrefs := captured[0].Tags()

	if len(gotPrefs) != 2 {
		t.Fatalf("expected 2 stored tags, got %d", len(gotPrefs))
	}
}

func TestUserService_SaveAndUnsaveContent(t *testing.T) {
	uid := uuid.New()
	cid := uuid.New()

	mockSaved := &mockSavedContentRepository{
		SaveFn: func(ctx context.Context, sc *domain.SavedContent) error {
			if sc.UserID() != uid || sc.SubjectID() != cid {
				return errors.New("wrong data passed to Save")
			}
			return nil
		},
		UnsaveFn: func(ctx context.Context, userID, contentID uuid.UUID, ct domain.ContentType) error {
			if userID != uid || contentID != cid || ct != domain.POST {
				return errors.New("wrong data passed to Unsave")
			}
			return nil
		},
	}

	svc := services.NewUserService(nil, nil, mockSaved)

	t.Run("save", func(t *testing.T) {
		err := svc.SaveContent(context.Background(), uid, cid, domain.POST)
		if err != nil {
			t.Fatalf("unexpected error in SaveContent: %v", err)
		}
	})

	t.Run("unsave", func(t *testing.T) {
		err := svc.UnsaveContent(context.Background(), uid, cid, domain.POST)
		if err != nil {
			t.Fatalf("unexpected error in UnsaveContent: %v", err)
		}
	})
}

func TestUserService_ListSavedContent(t *testing.T) {
	uid := uuid.New()
	page := filters.NewPage(10, 0)

	sc1, _ := domain.NewSavedContent(time.Now(), uid, uuid.New(), domain.POST)
	sc2, _ := domain.NewSavedContent(time.Now(), uid, uuid.New(), domain.COMMENT)

	mockSaved := &mockSavedContentRepository{
		ListFn: func(ctx context.Context, userID uuid.UUID, p filters.Page) ([]domain.SavedContent, error) {
			return []domain.SavedContent{*sc1, *sc2}, nil
		},
	}

	svc := services.NewUserService(nil, nil, mockSaved)

	out, err := svc.ListSavedContent(context.Background(), uid, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out))
	}
}
