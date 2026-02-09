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

func TestPostService_Create_Success(t *testing.T) {
	authorID := uuid.New()
	called := false

	mockRepo := &mockPostRepository{
		CreateFn: func(ctx context.Context, p *domain.Post) error {
			called = true
			if p.AuthorID() != authorID {
				return errors.New("wrong author id passed")
			}
			if p.Title() != "Hello" {
				return errors.New("wrong title")
			}
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	tags := []domain.Tag{}
	stances := make([]domain.Stance, 3)
	for i := range stances {
		color, _ := domain.NewColor("#ffffff")
		s, _ := domain.NewStance(uuid.New(), "stance", color, "")
		stances[i] = s
	}

	post, err := svc.Create(context.Background(), authorID, "Hello", "Desc", tags, stances)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post == nil {
		t.Fatalf("expected non-nil post")
	}
	if !called {
		t.Fatalf("expected repository.Create to be called")
	}
}

func TestPostService_Create_DomainError(t *testing.T) {
	mockRepo := &mockPostRepository{
		CreateFn: func(ctx context.Context, p *domain.Post) error {
			t.Fatalf("Create should not be called when domain.NewPost fails")
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	_, err := svc.Create(context.Background(), uuid.New(), "Title", "Desc", nil, nil)
	if err == nil {
		t.Fatalf("expected error from domain.NewPost, got nil")
	}
}

func TestPostService_Delete_Success(t *testing.T) {
	authorID := uuid.New()
	postID := uuid.New()

	post, _ := domain.NewPost(time.Now(), "T", "D", nil, makeValidStances(), authorID)

	mockRepo := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			if id != postID {
				return nil, domain.ErrNotFound
			}
			return post, nil
		},
		DeleteFn: func(ctx context.Context, id uuid.UUID) error {
			if id != postID {
				return errors.New("wrong id to Delete")
			}
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	err := svc.Delete(context.Background(), authorID, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostService_Delete_Forbidden(t *testing.T) {
	authorID := uuid.New()
	otherID := uuid.New()
	postID := uuid.New()

	post, _ := domain.NewPost(time.Now(), "T", "D", nil, makeValidStances(), authorID)

	mockRepo := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			return post, nil
		},
		DeleteFn: func(ctx context.Context, id uuid.UUID) error {
			t.Fatalf("Delete should not be called when forbidden")
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	err := svc.Delete(context.Background(), otherID, postID)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestPostService_Update_Success(t *testing.T) {
	authorID := uuid.New()
	postID := uuid.New()

	original, _ := domain.NewPost(time.Now().Add(-time.Hour), "Old", "OldDesc", nil, makeValidStances(), authorID)

	var updatedSaved *domain.Post

	mockRepo := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			if id != postID {
				return nil, domain.ErrNotFound
			}
			return original, nil
		},
		UpdateFn: func(ctx context.Context, p *domain.Post) error {
			updatedSaved = p
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	newStances := makeValidStances()
	newTags := []domain.Tag{}
	post, err := svc.Update(context.Background(), authorID, postID, "New", "NewDesc", newTags, newStances)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.Title() != "New" {
		t.Fatalf("expected title New, got %s", post.Title())
	}
	if post.Content() != "NewDesc" {
		t.Fatalf("expected desc NewDesc, got %s", post.Content())
	}
	if updatedSaved == nil {
		t.Fatalf("expected Update to be called with non-nil post")
	}
}

func TestPostService_Update_Forbidden(t *testing.T) {
	authorID := uuid.New()
	otherID := uuid.New()
	postID := uuid.New()

	original, _ := domain.NewPost(time.Now(), "T", "D", nil, makeValidStances(), authorID)

	mockRepo := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			return original, nil
		},
		UpdateFn: func(ctx context.Context, p *domain.Post) error {
			t.Fatalf("Update should not be called when forbidden")
			return nil
		},
	}

	svc := services.NewPostService(mockRepo)

	_, err := svc.Update(context.Background(), otherID, postID, "New", "NewDesc", nil, makeValidStances())
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestPostService_GetFeed(t *testing.T) {
	userID := uuid.New()
	page := filters.NewPage(10, 0)

	mockRepo := &mockPostRepository{
		GetFeedFn: func(ctx context.Context, uid uuid.UUID, p filters.Page) ([]domain.Post, error) {
			if uid != userID {
				return nil, errors.New("wrong user id")
			}
			return []domain.Post{}, nil
		},
	}

	svc := services.NewPostService(mockRepo)

	_, err := svc.GetFeed(context.Background(), userID, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostService_ListUserPosts(t *testing.T) {
	userID := uuid.New()
	page := filters.NewPage(5, 0)

	mockRepo := &mockPostRepository{
		ListByUserIDFn: func(ctx context.Context, uid uuid.UUID, p filters.Page) ([]domain.Post, error) {
			if uid != userID {
				return nil, errors.New("wrong user id")
			}
			return []domain.Post{}, nil
		},
	}

	svc := services.NewPostService(mockRepo)

	_, err := svc.ListUserPosts(context.Background(), userID, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	postID := uuid.New()

	post, _ := domain.NewPost(time.Now(), "T", "D", nil, makeValidStances(), uuid.New())

	mockRepo := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			if id != postID {
				return nil, domain.ErrNotFound
			}
			return post, nil
		},
	}

	svc := services.NewPostService(mockRepo)

	got, err := svc.GetPostByID(context.Background(), postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title() != "T" {
		t.Fatalf("expected title T, got %s", got.Title())
	}
}

func TestPostService_Search(t *testing.T) {
	mockRepo := &mockPostRepository{
		SearchFn: func(ctx context.Context, f filters.SearchPostFilter) ([]domain.Post, error) {
			return []domain.Post{}, nil
		},
	}

	svc := services.NewPostService(mockRepo)

	f := filters.SearchPostFilter{
		Page: filters.NewPage(10, 0),
	}
	uuid, _ := uuid.NewUUID()
	_, err := svc.Search(context.Background(), uuid, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func makeValidStances() []domain.Stance {
	st := make([]domain.Stance, 3)
	for i := range st {
		color, _ := domain.NewColor("#ffffff")
		s, _ := domain.NewStance(uuid.New(), "stance", color, "")
		st[i] = s
	}
	return st
}
