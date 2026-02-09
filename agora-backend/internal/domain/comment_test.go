package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

// ---- Tests for NewComment ----

func TestNewComment_Valid_NoReply(t *testing.T) {
	now := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	content := "  This is a comment  "

	c, err := domain.NewComment(now, content, stance, postID, authorID, nil)
	if err != nil {
		t.Fatalf("NewComment unexpected error: %v", err)
	}
	if c == nil {
		t.Fatalf("NewComment returned nil comment")
	}

	if c.ID() == uuid.Nil {
		t.Fatalf("expected non-nil comment ID")
	}
	if !c.CreatedAt().Equal(now) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", now, c.CreatedAt())
	}
	if !c.UpdatedAt().Equal(now) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", now, c.UpdatedAt())
	}

	if c.Content() == "" {
		t.Fatalf("expected non-empty content")
	}

	if c.Stance().ID() != stance.ID() {
		t.Fatalf("stance mismatch: expected ID %v, got %v", stance.ID(), c.Stance().ID())
	}

	if c.PostID() != postID {
		t.Fatalf("PostID mismatch: expected %v, got %v", postID, c.PostID())
	}
	if c.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, c.AuthorID())
	}
	if c.IsReply() {
		t.Fatalf("expected IsReply=false when replyToID is nil")
	}
	if c.ReplyToID() != nil {
		t.Fatalf("expected ReplyToID=nil")
	}
}

func TestNewComment_Valid_WithReply(t *testing.T) {
	now := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	replyTo := uuid.New()
	stance := mustStance(t, uuid.New(), "Disagree", "#ff0000", "")

	c, err := domain.NewComment(now, "content", stance, postID, authorID, &replyTo)
	if err != nil {
		t.Fatalf("NewComment unexpected error: %v", err)
	}

	if !c.IsReply() {
		t.Fatalf("expected IsReply=true when replyToID is non-nil")
	}
	if c.ReplyToID() == nil || *c.ReplyToID() != replyTo {
		t.Fatalf("ReplyToID mismatch: expected %v, got %v", replyTo, c.ReplyToID())
	}
}

func TestNewComment_EmptyContent(t *testing.T) {
	now := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	_, err := domain.NewComment(now, "   ", stance, postID, authorID, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrEmptyCommentContent {
		t.Fatalf("expected ErrEmptyCommentContent, got %v", err)
	}
}

func TestNewComment_ContentTooLong(t *testing.T) {
	now := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	longContent := strings.Repeat("a", 5001)

	_, err := domain.NewComment(now, longContent, stance, postID, authorID, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrCommentContentTooLong {
		t.Fatalf("expected ErrCommentContentTooLong, got %v", err)
	}
}

func TestNewComment_InvalidIDs(t *testing.T) {
	now := time.Now()
	validPostID := uuid.New()
	validAuthorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	replyID := uuid.New()

	tests := []struct {
		name     string
		postID   uuid.UUID
		authorID uuid.UUID
		replyPtr *uuid.UUID
	}{
		{"nil author", validPostID, uuid.Nil, &replyID},
		{"nil post", uuid.Nil, validAuthorID, &replyID},
		{"nil reply id pointer", validPostID, validAuthorID, &uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			print(tt.name)
			_, err := domain.NewComment(now, "content", stance, tt.postID, tt.authorID, tt.replyPtr)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidID {
				t.Fatalf("expected ErrInvalidID, got %v", err)
			}
		})
	}
}

// ---- Tests for ExistingComment ----

func TestExistingComment_Valid(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()

	postID := uuid.New()
	authorID := uuid.New()
	replyTo := uuid.New()
	stance := mustStance(t, uuid.New(), "Neutral", "#0000ff", "")

	c, err := domain.ExistingComment(id, created, updated, "content", stance, postID, authorID, &replyTo)
	if err != nil {
		t.Fatalf("ExistingComment unexpected error: %v", err)
	}

	if c.ID() != id {
		t.Fatalf("ID mismatch: expected %v, got %v", id, c.ID())
	}
	if !c.CreatedAt().Equal(created) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", created, c.CreatedAt())
	}
	if !c.UpdatedAt().Equal(updated) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", updated, c.UpdatedAt())
	}
	if c.PostID() != postID {
		t.Fatalf("PostID mismatch: expected %v, got %v", postID, c.PostID())
	}
	if c.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, c.AuthorID())
	}
	if !c.IsReply() || c.ReplyToID() == nil || *c.ReplyToID() != replyTo {
		t.Fatalf("ReplyToID mismatch: expected %v, got %v", replyTo, c.ReplyToID())
	}
}

func TestExistingComment_InvalidContent(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	_, err := domain.ExistingComment(id, created, updated, "   ", stance, postID, authorID, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrEmptyCommentContent {
		t.Fatalf("expected ErrEmptyCommentContent, got %v", err)
	}
}

func TestExistingComment_InvalidIDs(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	validPostID := uuid.New()
	validAuthorID := uuid.New()
	stance := mustStance(t, uuid.New(), "Agree", "#00ff00", "")

	nilUUID := uuid.Nil

	tests := []struct {
		name     string
		postID   uuid.UUID
		authorID uuid.UUID
		replyPtr *uuid.UUID
	}{
		{"nil author", validPostID, uuid.Nil, nil},
		{"nil post", uuid.Nil, validAuthorID, nil},
		{"nil reply id pointer", validPostID, validAuthorID, &nilUUID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ExistingComment(id, created, updated, "content", stance, tt.postID, tt.authorID, tt.replyPtr)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidID {
				t.Fatalf("expected ErrInvalidID, got %v", err)
			}
		})
	}
}
