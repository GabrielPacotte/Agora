package services

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/repositories"
	"github.com/google/uuid"
)

type PostService interface {
	Create(ctx context.Context, requesterID uuid.UUID, title, description string, tags []domain.Tag, stances []domain.Stance) (*domain.Post, error)
	Delete(ctx context.Context, requesterID, postID uuid.UUID) error
	Update(ctx context.Context, requesterID, postID uuid.UUID, newTitle, newDesc string, newTags []domain.Tag, newStances []domain.Stance) (*domain.Post, error)

	GetFeed(ctx context.Context, requesterID uuid.UUID, page filters.Page) ([]domain.Post, error)
	ListUserPosts(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error)
	Search(ctx context.Context, requesterID uuid.UUID, searchFilter filters.SearchPostFilter) ([]domain.Post, error)
}

type postService struct {
	posts repositories.PostRepository
}

func NewPostService(posts repositories.PostRepository) PostService {
	return &postService{posts: posts}
}

func (s *postService) Create(
	ctx context.Context,
	requesterID uuid.UUID,
	title, description string,
	tags []domain.Tag,
	stances []domain.Stance,
) (*domain.Post, error) {
	now := time.Now().UTC()

	post, err := domain.NewPost(now, title, description, tags, stances, requesterID)
	if err != nil {
		return nil, err
	}

	if err := s.posts.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) Delete(ctx context.Context, requesterID, postID uuid.UUID) error {
	post, err := s.posts.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID() != requesterID {
		return domain.ErrForbidden
	}

	return s.posts.Delete(ctx, postID)
}

func (s *postService) Update(
	ctx context.Context,
	requesterID, postID uuid.UUID,
	newTitle, newDesc string,
	newTags []domain.Tag,
	newStances []domain.Stance,
) (*domain.Post, error) {
	existing, err := s.posts.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID() != requesterID {
		return nil, domain.ErrForbidden
	}

	now := time.Now().UTC()

	updated, err := domain.ExistingPost(
		existing.ID(),
		existing.CreatedAt(),
		now,
		newTitle,
		newDesc,
		newTags,
		newStances,
		existing.AuthorID(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.posts.Update(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *postService) GetFeed(ctx context.Context, requesterID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	if requesterID == uuid.Nil {
		return nil, domain.ErrForbidden
	}

	page = filters.NewPage(page.Limit(), page.Offset())
	return s.posts.GetFeed(ctx, requesterID, page)
}

func (s *postService) Search(ctx context.Context, requesterID uuid.UUID, f filters.SearchPostFilter) ([]domain.Post, error) {
	if requesterID == uuid.Nil {
		return nil, domain.ErrForbidden
	}

	f = normalizeSearchPostFilter(f)
	return s.posts.Search(ctx, f)
}

func normalizeSearchPostFilter(f filters.SearchPostFilter) filters.SearchPostFilter {
	f.Page = filters.NewPage(f.Limit(), f.Offset())
	if f.FromDate != nil && f.ToDate != nil && f.ToDate.Before(*f.FromDate) {
		f.ToDate = nil
	}
	return f
}

func (s *postService) ListUserPosts(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	return s.posts.ListByUserID(ctx, userID, page)
}

func (s *postService) GetPostByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error) {
	return s.posts.GetByID(ctx, postID)
}
