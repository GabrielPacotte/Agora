package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

// Helpers

func mustColor(t *testing.T, hex string) domain.Color {
	t.Helper()
	c, err := domain.NewColor(hex)
	if err != nil {
		t.Fatalf("NewColor(%q) unexpected error: %v", hex, err)
	}
	return c
}

func mustTag(t *testing.T, v string) domain.Tag {
	t.Helper()
	tag, err := domain.NewTag(v)
	if err != nil {
		t.Fatalf("NewTag(%q) unexpected error: %v", v, err)
	}
	return tag
}

func mustStance(t *testing.T, id uuid.UUID, label, hexColor, desc string) domain.Stance {
	t.Helper()
	color, err := domain.NewColor(hexColor)
	if err != nil {
		t.Fatalf("NewColor(%q) unexpected error: %v", hexColor, err)
	}
	s, err := domain.NewStance(id, label, color, desc)
	if err != nil {
		t.Fatalf("NewStance(%q) unexpected error: %v", label, err)
	}
	return s
}

// ---- Tests for NewPost ----

func TestNewPost_Valid(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	tag1 := mustTag(t, "politics")
	tag2 := mustTag(t, "france")

	s1 := mustStance(t, uuid.New(), "Agree", "#00ff00", "")
	s2 := mustStance(t, uuid.New(), "Disagree", "#ff0000", "")
	s3 := mustStance(t, uuid.New(), "Neutral", "#0000ff", "")

	title := "   My title   "
	desc := "   Some Content   "

	p, err := domain.NewPost(now, title, desc, []domain.Tag{tag1, tag2}, []domain.Stance{s1, s2, s3}, authorID)
	if err != nil {
		t.Fatalf("NewPost unexpected error: %v", err)
	}
	if p == nil {
		t.Fatalf("NewPost returned nil post")
	}

	// Title & Content must be trimmed
	if p.Title() != strings.TrimSpace(title) {
		t.Fatalf("Title mismatch: expected %q, got %q", strings.TrimSpace(title), p.Title())
	}
	if p.Content() != strings.TrimSpace(desc) {
		t.Fatalf("Content mismatch: expected %q, got %q", strings.TrimSpace(desc), p.Content())
	}
	if !p.HasContent() {
		t.Fatalf("HasContent should be true when Content is non-empty")
	}

	// Author
	if p.AuthorID() != authorID {
		t.Fatalf("AuthorID mismatch: expected %v, got %v", authorID, p.AuthorID())
	}

	// Tags
	gotTags := p.Tags()
	if len(gotTags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(gotTags))
	}
	if gotTags[0] != tag1 || gotTags[1] != tag2 {
		t.Fatalf("tags mismatch: expected [%v %v], got [%v %v]", tag1, tag2, gotTags[0], gotTags[1])
	}

	// Stances
	gotStances := p.Stances()
	if len(gotStances) != 3 {
		t.Fatalf("expected 3 stances, got %d", len(gotStances))
	}

	// ID / timestamps via embedded identified
	if p.ID() == uuid.Nil {
		t.Fatalf("expected non-nil post ID")
	}
	if !p.CreatedAt().Equal(now) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", now, p.CreatedAt())
	}
	if !p.UpdatedAt().Equal(now) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", now, p.UpdatedAt())
	}
}

func TestNewPost_EmptyTitle(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	_, err := domain.NewPost(now, "   ", "desc", nil, stances, authorID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrEmptyPostTitle) {
		t.Fatalf("expected ErrEmptyPostTitle, got %v", err)
	}
}

func TestNewPost_TitleTooLong(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	longTitle := strings.Repeat("a", domain.MaxPostTitleLength+1)

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	_, err := domain.NewPost(now, longTitle, "desc", nil, stances, authorID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrPostTitleTooLong) {
		t.Fatalf("expected ErrPostTitleTooLong, got %v", err)
	}
}

func TestNewPost_ContentTooLong(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()
	longDesc := strings.Repeat("d", domain.MaxPostContentLength+1)

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	_, err := domain.NewPost(now, "Title", longDesc, nil, stances, authorID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrPostContentTooLong) {
		t.Fatalf("expected ErrPostContentTooLong, got %v", err)
	}
}

func TestNewPost_NotEnoughStances(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
	}

	_, err := domain.NewPost(now, "Title", "desc", nil, stances, authorID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrMinimalStanceAmount) {
		t.Fatalf("expected ErrMinimalStanceAmount, got %v", err)
	}
}

func TestNewPost_InvalidAuthorID(t *testing.T) {
	now := time.Now()
	authorID := uuid.Nil

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	_, err := domain.NewPost(now, "Title", "desc", nil, stances, authorID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestNewPost_EmptyContentHasContentFalse(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	p, err := domain.NewPost(now, "Title", "   ", nil, stances, authorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Content() != "" {
		t.Fatalf("expected empty Content after trim, got %q", p.Content())
	}
	if p.HasContent() {
		t.Fatalf("expected HasContent=false for empty Content")
	}
}

// ---- Tests for dedupe behavior ----

func TestNewPost_DedupeTagsAndStances(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	tagA := mustTag(t, "tag-a")
	tagB := mustTag(t, "tag-b")

	tags := []domain.Tag{tagA, tagB, tagA, tagB}

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	s1 := mustStance(t, id1, "A", "#111111", "")
	s2 := mustStance(t, id2, "B", "#222222", "")
	s1dup := mustStance(t, id1, "A duplicate", "#111111", "whatever")
	s3 := mustStance(t, id3, "C", "#333333", "")

	stances := []domain.Stance{s1, s2, s1dup, s3}
	p, err := domain.NewPost(now, "Title", "Desc", tags, stances, authorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotTags := p.Tags()
	if len(gotTags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d", len(gotTags))
	}
	if gotTags[0] != tagA || gotTags[1] != tagB {
		t.Fatalf("unexpected tags order/content: %#v", gotTags)
	}

	gotStances := p.Stances()
	if len(gotStances) != 3 {
		t.Fatalf("expected 3 unique stances, got %d", len(gotStances))
	}
	if gotStances[0].ID() != id1 || gotStances[1].ID() != id2 {
		t.Fatalf("unexpected stances order/content: [%v, %v]", gotStances[0].ID(), gotStances[1].ID())
	}
}

// ---- Tests for copies in getters ----

func TestPost_TagsReturnsCopy(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	tagA := mustTag(t, "tag-a")
	tagB := mustTag(t, "tag-b")

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	p, err := domain.NewPost(now, "Title", "Desc", []domain.Tag{tagA, tagB}, stances, authorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tags1 := p.Tags()
	tags1[0] = mustTag(t, "hacked")

	tags2 := p.Tags()
	if tags2[0] != tagA {
		t.Fatalf("expected original tag at index 0 to remain %q, got %q", tagA, tags2[0])
	}
}

func TestPost_StancesReturnsCopy(t *testing.T) {
	now := time.Now()
	authorID := uuid.New()

	s1 := mustStance(t, uuid.New(), "A", "#000001", "")
	s2 := mustStance(t, uuid.New(), "B", "#000002", "")
	s3 := mustStance(t, uuid.New(), "C", "#000003", "")

	p, err := domain.NewPost(now, "Title", "Desc", nil, []domain.Stance{s1, s2, s3}, authorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	st1 := p.Stances()
	st1[0] = mustStance(t, uuid.New(), "Hacked", "#ffffff", "")

	st2 := p.Stances()
	if st2[0].ID() != s1.ID() {
		t.Fatalf("expected original stance at index 0 to remain, got different ID")
	}
}

// ---- Tests for ExistingPost ----

func TestExistingPost_Valid(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	authorID := uuid.New()

	tag := mustTag(t, "tag")
	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	p, err := domain.ExistingPost(id, created, updated, "  Title  ", "  Desc  ", []domain.Tag{tag}, stances, authorID)
	if err != nil {
		t.Fatalf("ExistingPost unexpected error: %v", err)
	}

	if p.ID() != id {
		t.Fatalf("ID mismatch: expected %v, got %v", id, p.ID())
	}
	if !p.CreatedAt().Equal(created) {
		t.Fatalf("CreatedAt mismatch: expected %v, got %v", created, p.CreatedAt())
	}
	if !p.UpdatedAt().Equal(updated) {
		t.Fatalf("UpdatedAt mismatch: expected %v, got %v", updated, p.UpdatedAt())
	}
}

func TestExistingPost_InvalidIDOrTimestamps(t *testing.T) {
	id := uuid.New()
	created := time.Now().Add(-time.Hour)
	updated := time.Now()
	authorID := uuid.New()

	stances := []domain.Stance{
		mustStance(t, uuid.New(), "A", "#000001", ""),
		mustStance(t, uuid.New(), "B", "#000002", ""),
		mustStance(t, uuid.New(), "C", "#000003", ""),
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		created time.Time
		updated time.Time
		wantErr error
	}{
		{"nil ID", uuid.Nil, created, updated, domain.ErrInvalidID},
		{"createdAt zero", id, time.Time{}, updated, domain.ErrInvalidTimestamp},
		{"updatedAt zero", id, created, time.Time{}, domain.ErrInvalidTimestamp},
		{"updated before created", id, updated, created, domain.ErrInvalidTimestamp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ExistingPost(tt.id, tt.created, tt.updated, "Title", "Desc", nil, stances, authorID)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
