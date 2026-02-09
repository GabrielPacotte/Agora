package services

import (
	"context"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	Create(ctx context.Context, requesterID, postID uuid.UUID, stanceID uuid.UUID, content string, replyToID *uuid.UUID) (*domain.Comment, error)
	Update(ctx context.Context, requesterID, commentID uuid.UUID, newContent string) (*domain.Comment, error)
	Delete(ctx context.Context, requesterID, commentID uuid.UUID) error

	ListRootComments(ctx context.Context, postID uuid.UUID, page filters.Page) ([]domain.Comment, error)
	ListReplies(ctx context.Context, parentCommentID uuid.UUID, page filters.Page) ([]domain.Comment, error)
}

type commentService struct {
	comments repositories.CommentRepository
	posts    repositories.PostRepository
}

func NewCommentService(
	comments repositories.CommentRepository,
	posts repositories.PostRepository,
) CommentService {
	return &commentService{
		comments: comments,
		posts:    posts,
	}
}

func (s *commentService) Create(
	ctx context.Context,
	requesterID, postID uuid.UUID,
	stanceID uuid.UUID,
	content string,
	replyToID *uuid.UUID,
) (*domain.Comment, error) {
	post, err := s.posts.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	var stance domain.Stance
	found := false
	for _, st := range post.Stances() {
		if st.ID() == stanceID {
			stance = st
			found = true
			break
		}
	}
	if !found {
		return nil, domain.ErrNotFound
	}

	if replyToID != nil {
		replyToComment, err := s.comments.GetByID(ctx, *replyToID)
		if err != nil {
			return nil, err
		}
		if replyToComment.PostID() != postID {
			return nil, domain.ErrNotFound
		}
	}

	now := time.Now().UTC()

	comment, err := domain.NewComment(now, content, stance, postID, requesterID, replyToID)
	if err != nil {
		return nil, err
	}

	if err := s.comments.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) Update(
	ctx context.Context,
	requesterID, commentID uuid.UUID,
	newContent string,
) (*domain.Comment, error) {
	existing, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID() != requesterID {
		return nil, domain.ErrForbidden
	}

	now := time.Now().UTC()

	updated, err := domain.ExistingComment(
		existing.ID(),
		existing.CreatedAt(),
		now,
		newContent,
		existing.Stance(),
		existing.PostID(),
		existing.AuthorID(),
		existing.ReplyToID(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.comments.Update(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *commentService) Delete(
	ctx context.Context,
	requesterID, commentID uuid.UUID,
) error {
	existing, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if existing.AuthorID() != requesterID {
		return domain.ErrForbidden
	}

	return s.comments.Delete(ctx, commentID)
}

func (s *commentService) ListRootComments(
	ctx context.Context,
	postID uuid.UUID,
	page filters.Page,
) ([]domain.Comment, error) {
	return s.comments.ListRootComments(ctx, postID, page)
}

func (s *commentService) ListReplies(
	ctx context.Context,
	parentCommentID uuid.UUID,
	page filters.Page,
) ([]domain.Comment, error) {
	return s.comments.ListReplies(ctx, parentCommentID, page)
}
