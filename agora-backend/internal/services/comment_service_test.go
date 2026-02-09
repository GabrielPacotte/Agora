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

func TestCommentService_Create_Success(t *testing.T) {
	ctx := context.Background()

	requesterID := uuid.New()
	postID := uuid.New()

	stances := makeValidStances()
	chosenStance := stances[0]

	post, err := domain.NewPost(time.Now(), "Title", "Desc", nil, stances, requesterID)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}

	mockPosts := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			if id != postID {
				return nil, domain.ErrNotFound
			}
			return post, nil
		},
	}

	var createdComment *domain.Comment
	mockComments := &mockCommentRepository{
		CreateFn: func(ctx context.Context, c *domain.Comment) error {
			createdComment = c
			if c.AuthorID() != requesterID {
				return errors.New("wrong author id")
			}
			if c.PostID() != postID {
				return errors.New("wrong post id")
			}
			if c.Stance().ID() != chosenStance.ID() {
				return errors.New("wrong stance id")
			}
			if c.Content() != "Hello world" {
				return errors.New("wrong content")
			}
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, mockPosts)

	comment, err := svc.Create(ctx, requesterID, postID, chosenStance.ID(), "Hello world", nil)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if comment == nil {
		t.Fatalf("expected non-nil comment")
	}
	if createdComment == nil {
		t.Fatalf("expected repository.Create to be called")
	}
}

func TestCommentService_Create_UnknownStance(t *testing.T) {
	ctx := context.Background()

	requesterID := uuid.New()
	postID := uuid.New()

	stances := makeValidStances()
	post, _ := domain.NewPost(time.Now(), "Title", "Desc", nil, stances, requesterID)

	mockPosts := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			return post, nil
		},
	}

	mockComments := &mockCommentRepository{
		CreateFn: func(ctx context.Context, c *domain.Comment) error {
			t.Fatalf("Create should not be called when stance is unknown")
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, mockPosts)

	unknownStanceID := uuid.New()

	_, err := svc.Create(ctx, requesterID, postID, unknownStanceID, "Hello", nil)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown stance, got %v", err)
	}
}

func TestCommentService_Create_PostNotFound(t *testing.T) {
	ctx := context.Background()

	mockPosts := &mockPostRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
			return nil, domain.ErrNotFound
		},
	}

	mockComments := &mockCommentRepository{
		CreateFn: func(ctx context.Context, c *domain.Comment) error {
			t.Fatalf("Create should not be called when post not found")
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, mockPosts)

	_, err := svc.Create(ctx, uuid.New(), uuid.New(), uuid.New(), "Hello", nil)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCommentService_Update_Success(t *testing.T) {
	ctx := context.Background()

	authorID := uuid.New()
	commentID := uuid.New()
	postID := uuid.New()

	stances := makeValidStances()
	stance := stances[0]

	existing, err := domain.NewComment(
		time.Now().Add(-time.Hour),
		"Old content",
		stance,
		postID,
		authorID,
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment error: %v", err)
	}

	var updatedSaved *domain.Comment

	mockComments := &mockCommentRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
			if id != commentID {
				return nil, domain.ErrNotFound
			}
			return existing, nil
		},
		UpdateFn: func(ctx context.Context, c *domain.Comment) error {
			updatedSaved = c
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	newContent := "Updated content"
	updated, err := svc.Update(ctx, authorID, commentID, newContent)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if updated.Content() != newContent {
		t.Fatalf("expected content %q, got %q", newContent, updated.Content())
	}
	if updatedSaved == nil {
		t.Fatalf("expected repository.Update to be called")
	}
}

func TestCommentService_Update_Forbidden(t *testing.T) {
	ctx := context.Background()

	authorID := uuid.New()
	otherID := uuid.New()
	commentID := uuid.New()
	postID := uuid.New()

	stance := makeValidStances()[0]

	existing, err := domain.NewComment(
		time.Now(),
		"Old",
		stance,
		postID,
		authorID,
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment error: %v", err)
	}

	mockComments := &mockCommentRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
			return existing, nil
		},
		UpdateFn: func(ctx context.Context, c *domain.Comment) error {
			t.Fatalf("Update should not be called when forbidden")
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	_, err = svc.Update(ctx, otherID, commentID, "New")
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCommentService_Delete_Success(t *testing.T) {
	ctx := context.Background()

	authorID := uuid.New()
	commentID := uuid.New()
	postID := uuid.New()
	stance := makeValidStances()[0]

	existing, err := domain.NewComment(
		time.Now(),
		"Some content",
		stance,
		postID,
		authorID,
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment error: %v", err)
	}

	mockComments := &mockCommentRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
			if id != commentID {
				return nil, domain.ErrNotFound
			}
			return existing, nil
		},
		DeleteFn: func(ctx context.Context, id uuid.UUID) error {
			if id != commentID {
				return errors.New("wrong comment id")
			}
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	if err := svc.Delete(ctx, authorID, commentID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestCommentService_Delete_Forbidden(t *testing.T) {
	ctx := context.Background()

	authorID := uuid.New()
	otherID := uuid.New()
	commentID := uuid.New()
	postID := uuid.New()
	stance := makeValidStances()[0]

	existing, err := domain.NewComment(
		time.Now(),
		"Some content",
		stance,
		postID,
		authorID,
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment error: %v", err)
	}

	mockComments := &mockCommentRepository{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
			return existing, nil
		},
		DeleteFn: func(ctx context.Context, id uuid.UUID) error {
			t.Fatalf("Delete should not be called when forbidden")
			return nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	err = svc.Delete(ctx, otherID, commentID)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCommentService_ListRootComments(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	page := filters.NewPage(10, 0)

	mockComments := &mockCommentRepository{
		ListRootCommentsFn: func(ctx context.Context, pid uuid.UUID, p filters.Page) ([]domain.Comment, error) {
			if pid != postID {
				return nil, errors.New("wrong post id")
			}
			return []domain.Comment{}, nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	_, err := svc.ListRootComments(ctx, postID, page)
	if err != nil {
		t.Fatalf("ListRootComments error: %v", err)
	}
}

func TestCommentService_ListReplies(t *testing.T) {
	ctx := context.Background()
	commentID := uuid.New()
	page := filters.NewPage(10, 0)

	mockComments := &mockCommentRepository{
		ListRepliesFn: func(ctx context.Context, cid uuid.UUID, p filters.Page) ([]domain.Comment, error) {
			if cid != commentID {
				return nil, errors.New("wrong comment id")
			}
			return []domain.Comment{}, nil
		},
	}

	svc := services.NewCommentService(mockComments, nil)

	_, err := svc.ListReplies(ctx, commentID, page)
	if err != nil {
		t.Fatalf("ListReplies error: %v", err)
	}
}
