package services_test

import (
	"context"
	"errors"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/google/uuid"
)

type mockUserRepository struct {
	GetByIDFn                    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByLoginWithPasswordHashFn func(ctx context.Context, login string) (*domain.User, string, error)
	DeleteFn                     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

func (m *mockUserRepository) Create(context.Context, *domain.User, string) error {
	return errors.New("not implemented")
}
func (m *mockUserRepository) UpdatePublicFields(context.Context, *domain.User) error {
	return errors.New("not implemented")
}
func (m *mockUserRepository) UpdatePassword(context.Context, *domain.User, string) error {
	return errors.New("not implemented")
}

type mockUserPreferencesRepository struct {
	ListFn   func(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error)
	UpdateFn func(ctx context.Context, prefs *domain.UserPreferences) error
}

func (m *mockUserPreferencesRepository) List(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error) {
	return m.ListFn(ctx, userID)
}

func (m *mockUserPreferencesRepository) Update(ctx context.Context, prefs *domain.UserPreferences) error {
	return m.UpdateFn(ctx, prefs)
}

type mockSavedContentRepository struct {
	SaveFn   func(ctx context.Context, sc *domain.SavedContent) error
	UnsaveFn func(ctx context.Context, userID, contentID uuid.UUID, ct domain.ContentType) error
	ListFn   func(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.SavedContent, error)
}

func (m *mockSavedContentRepository) Save(ctx context.Context, sc *domain.SavedContent) error {
	return m.SaveFn(ctx, sc)
}

func (m *mockSavedContentRepository) Unsave(ctx context.Context, userID, contentID uuid.UUID, ct domain.ContentType) error {
	return m.UnsaveFn(ctx, userID, contentID, ct)
}

func (m *mockSavedContentRepository) List(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.SavedContent, error) {
	return m.ListFn(ctx, userID, page)
}

type mockPostRepository struct {
	CreateFn       func(ctx context.Context, post *domain.Post) error
	UpdateFn       func(ctx context.Context, post *domain.Post) error
	DeleteFn       func(ctx context.Context, postID uuid.UUID) error
	ListByUserIDFn func(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error)
	GetFeedFn      func(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error)
	GetByIDFn      func(ctx context.Context, postID uuid.UUID) (*domain.Post, error)
	SearchFn       func(ctx context.Context, f filters.SearchPostFilter) ([]domain.Post, error)
}

func (m *mockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	return m.CreateFn(ctx, post)
}

func (m *mockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	return m.UpdateFn(ctx, post)
}

func (m *mockPostRepository) Delete(ctx context.Context, postID uuid.UUID) error {
	return m.DeleteFn(ctx, postID)
}

func (m *mockPostRepository) ListByUserID(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	return m.ListByUserIDFn(ctx, userID, page)
}

func (m *mockPostRepository) GetFeed(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	return m.GetFeedFn(ctx, userID, page)
}

func (m *mockPostRepository) GetByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error) {
	return m.GetByIDFn(ctx, postID)
}

func (m *mockPostRepository) Search(ctx context.Context, f filters.SearchPostFilter) ([]domain.Post, error) {
	return m.SearchFn(ctx, f)
}

type mockCommentRepository struct {
	CreateFn           func(ctx context.Context, comment *domain.Comment) error
	UpdateFn           func(ctx context.Context, comment *domain.Comment) error
	DeleteFn           func(ctx context.Context, commentID uuid.UUID) error
	ListRootCommentsFn func(ctx context.Context, postID uuid.UUID, page filters.Page) ([]domain.Comment, error)
	ListRepliesFn      func(ctx context.Context, parentCommentID uuid.UUID, page filters.Page) ([]domain.Comment, error)
	GetByIDFn          func(ctx context.Context, commentID uuid.UUID) (*domain.Comment, error)
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return m.CreateFn(ctx, comment)
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return m.UpdateFn(ctx, comment)
}

func (m *mockCommentRepository) Delete(ctx context.Context, commentID uuid.UUID) error {
	return m.DeleteFn(ctx, commentID)
}

func (m *mockCommentRepository) ListRootComments(ctx context.Context, postID uuid.UUID, page filters.Page) ([]domain.Comment, error) {
	return m.ListRootCommentsFn(ctx, postID, page)
}

func (m *mockCommentRepository) ListReplies(ctx context.Context, parentCommentID uuid.UUID, page filters.Page) ([]domain.Comment, error) {
	return m.ListRepliesFn(ctx, parentCommentID, page)
}

func (m *mockUserRepository) GetByLoginWithPasswordHash(
	ctx context.Context,
	login string,
) (*domain.User, string, error) {
	if m.GetByLoginWithPasswordHashFn == nil {
		panic("GetByLoginWithPasswordHashFn not implemented")
	}
	return m.GetByLoginWithPasswordHashFn(ctx, login)
}

func (m *mockCommentRepository) GetByID(ctx context.Context, commentID uuid.UUID) (*domain.Comment, error) {
	return m.GetByIDFn(ctx, commentID)
}
