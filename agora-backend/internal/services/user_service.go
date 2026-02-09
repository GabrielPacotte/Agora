package services

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	DeleteAccount(ctx context.Context, userID uuid.UUID) error

	ListPreferences(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error)
	SetPreferences(ctx context.Context, userID uuid.UUID, preferences []domain.Tag) error

	SaveContent(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, contentType domain.ContentType) error
	UnsaveContent(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, contentType domain.ContentType) error
	ListSavedContent(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.SavedContent, error)
}

type userService struct {
	users        repositories.UserRepository
	preferences  repositories.UserPreferencesRepository
	savedContent repositories.SavedContentRepository
}

func NewUserService(
	users repositories.UserRepository,
	preferences repositories.UserPreferencesRepository,
	savedContent repositories.SavedContentRepository,
) UserService {
	return &userService{
		users:        users,
		preferences:  preferences,
		savedContent: savedContent,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *userService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return s.users.Delete(ctx, userID)
}

func (s *userService) ListPreferences(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	prefs, err := s.preferences.List(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return []domain.Tag{}, nil
		}
		return nil, err
	}
	tags := prefs.Tags()
	out := make([]domain.Tag, len(tags))
	copy(out, tags)
	return out, nil
}

func (s *userService) SetPreferences(ctx context.Context, userID uuid.UUID, preferences []domain.Tag) error {
	prefs, err := domain.NewUserPreferences(userID, preferences)
	if err != nil {
		return err
	}
	return s.preferences.Update(ctx, prefs)
}

func (s *userService) SaveContent(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, contentType domain.ContentType) error {
	now := time.Now().UTC()
	sc, err := domain.NewSavedContent(now, userID, contentID, contentType)
	if err != nil {
		return err
	}
	return s.savedContent.Save(ctx, sc)
}

func (s *userService) UnsaveContent(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, contentType domain.ContentType) error {
	return s.savedContent.Unsave(ctx, userID, contentID, contentType)
}

func (s *userService) ListSavedContent(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.SavedContent, error) {
	return s.savedContent.List(ctx, userID, page)
}
