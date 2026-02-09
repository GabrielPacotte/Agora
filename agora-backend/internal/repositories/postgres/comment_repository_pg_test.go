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

func TestCommentRepositoryPG_CreateAndListRootComments(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	commentRepo := postgres.NewCommentRepositoryPG(testPool)
	postRepo := postgres.NewPostRepositoryPG(testPool)
	userRepo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	author, err := domain.NewUser(now, "comment_author", "Comment Author")
	if err != nil {
		t.Fatalf("NewUser error: %v", err)
	}
	if err := userRepo.Create(ctx, author, "hash123"); err != nil {
		t.Fatalf("userRepo.Create error: %v", err)
	}

	color1, _ := domain.NewColor("#ff0000")
	color2, _ := domain.NewColor("#00ff00")
	color3, _ := domain.NewColor("#0000ff")

	s1, _ := domain.NewStance(uuid.New(), "Agree", color1, "Agree desc")
	s2, _ := domain.NewStance(uuid.New(), "Disagree", color2, "Disagree desc")
	s3, _ := domain.NewStance(uuid.New(), "Neutral", color3, "")

	post, err := domain.NewPost(
		now,
		"My Post for comments",
		"Body",
		nil,
		[]domain.Stance{s1, s2, s3},
		author.ID(),
	)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("postRepo.Create error: %v", err)
	}

	comment, err := domain.NewComment(
		now,
		"First root comment",
		s1,
		post.ID(),
		author.ID(),
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment error: %v", err)
	}

	if err := commentRepo.Create(ctx, comment); err != nil {
		t.Fatalf("commentRepo.Create error: %v", err)
	}

	page := filters.NewPage(10, 0)
	comments, err := commentRepo.ListRootComments(ctx, post.ID(), page)
	if err != nil {
		t.Fatalf("ListRootComments error: %v", err)
	}

	if len(comments) != 1 {
		t.Fatalf("expected 1 root comment, got %d", len(comments))
	}

	got := comments[0]

	if got.ID() != comment.ID() {
		t.Errorf("ID mismatch: got %v, want %v", got.ID(), comment.ID())
	}
	if got.PostID() != post.ID() {
		t.Errorf("PostID mismatch: got %v, want %v", got.PostID(), post.ID())
	}
	if got.AuthorID() != author.ID() {
		t.Errorf("AuthorID mismatch: got %v, want %v", got.AuthorID(), author.ID())
	}
	if got.Content() != comment.Content() {
		t.Errorf("Content mismatch: got %q, want %q", got.Content(), comment.Content())
	}
	if got.IsReply() {
		t.Errorf("expected root comment IsReply=false, got true")
	}

	if got.Stance().ID() != s1.ID() {
		t.Errorf("Stance ID mismatch: got %v, want %v", got.Stance().ID(), s1.ID())
	}
	if got.Stance().Label() != s1.Label() {
		t.Errorf("Stance Label mismatch: got %q, want %q", got.Stance().Label(), s1.Label())
	}
}

func TestCommentRepositoryPG_CreateReplyAndListReplies(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	commentRepo := postgres.NewCommentRepositoryPG(testPool)
	postRepo := postgres.NewPostRepositoryPG(testPool)
	userRepo := postgres.NewUserRepositoryPG(testPool)

	now := time.Now().UTC()

	author, err := domain.NewUser(now, "reply_author", "Reply Author")
	if err != nil {
		t.Fatalf("NewUser error: %v", err)
	}
	if err := userRepo.Create(ctx, author, "hash"); err != nil {
		t.Fatalf("userRepo.Create error: %v", err)
	}

	color, _ := domain.NewColor("#aaaaaa")
	s1, _ := domain.NewStance(uuid.New(), "Opinion", color, "")
	s2, _ := domain.NewStance(uuid.New(), "Other", color, "")
	s3, _ := domain.NewStance(uuid.New(), "Something", color, "")

	post, err := domain.NewPost(
		now,
		"Post with replies",
		"Body",
		nil,
		[]domain.Stance{s1, s2, s3},
		author.ID(),
	)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("postRepo.Create error: %v", err)
	}

	parent, err := domain.NewComment(
		now,
		"Parent comment",
		s1,
		post.ID(),
		author.ID(),
		nil,
	)
	if err != nil {
		t.Fatalf("NewComment(parent) error: %v", err)
	}
	if err := commentRepo.Create(ctx, parent); err != nil {
		t.Fatalf("commentRepo.Create(parent) error: %v", err)
	}

	replyToID := parent.ID()
	reply, err := domain.NewComment(
		now,
		"Child reply",
		s1,
		post.ID(),
		author.ID(),
		&replyToID,
	)
	if err != nil {
		t.Fatalf("NewComment(reply) error: %v", err)
	}
	if err := commentRepo.Create(ctx, reply); err != nil {
		t.Fatalf("commentRepo.Create(reply) error: %v", err)
	}

	page := filters.NewPage(10, 0)
	replies, err := commentRepo.ListReplies(ctx, parent.ID(), page)
	if err != nil {
		t.Fatalf("ListReplies error: %v", err)
	}

	if len(replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(replies))
	}

	got := replies[0]

	if !got.IsReply() {
		t.Errorf("expected IsReply=true, got false")
	}
	if got.ReplyToID() == nil || *got.ReplyToID() != parent.ID() {
		t.Errorf("ReplyToID mismatch: got %v, want %v", got.ReplyToID(), parent.ID())
	}
	if got.Content() != reply.Content() {
		t.Errorf("Content mismatch: got %q, want %q", got.Content(), reply.Content())
	}
	if got.Stance().ID() != s1.ID() {
		t.Errorf("Stance ID mismatch: got %v, want %v", got.Stance().ID(), s1.ID())
	}
}

func TestCommentRepositoryPG_Delete_NotFound(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	commentRepo := postgres.NewCommentRepositoryPG(testPool)

	randomID := uuid.New()

	err := commentRepo.Delete(ctx, randomID)
	if err == nil {
		t.Fatalf("expected error on deleting non-existing comment, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
