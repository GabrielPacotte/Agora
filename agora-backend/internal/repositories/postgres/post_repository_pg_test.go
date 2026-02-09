package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	"github.com/GabrielPacotte/Agora/internal/repositories/postgres"
	"github.com/google/uuid"
)

func TestPostRepositoryPG_CreateAndGetByID(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	repo := postgres.NewPostRepositoryPG(testPool)

	now := time.Now().UTC()
	author, err := domain.NewUser(now, "postauthor", "Post Author")
	if err != nil {
		t.Fatalf("NewUser error: %v", err)
	}

	userRepo := postgres.NewUserRepositoryPG(testPool)
	if err := userRepo.Create(ctx, author, "dummyHash"); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	tags := []domain.Tag{
		domain.Tag("politics"),
		domain.Tag("france"),
	}

	color1, _ := domain.NewColor("#ff0000")
	color2, _ := domain.NewColor("#00ff00")
	color3, _ := domain.NewColor("#0000ff")

	st1, _ := domain.NewStance(uuid.New(), "Pro", color1, "I support this")
	st2, _ := domain.NewStance(uuid.New(), "Against", color2, "I oppose this")
	st3, _ := domain.NewStance(uuid.New(), "Neutral", color3, "")

	post, err := domain.NewPost(
		now,
		"My Title",
		"My content",
		tags,
		[]domain.Stance{st1, st2, st3},
		author.ID(),
	)
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}

	if err := repo.Create(ctx, post); err != nil {
		t.Fatalf("repo.Create error: %v", err)
	}

	got, err := repo.GetByID(ctx, post.ID())
	if err != nil {
		t.Fatalf("repo.GetByID error: %v", err)
	}
	if got == nil {
		t.Fatalf("repo.GetByID returned nil")
	}

	if got.ID() != post.ID() {
		t.Errorf("ID mismatch: got %v want %v", got.ID(), post.ID())
	}

	if got.Title() != post.Title() {
		t.Errorf("Title mismatch: got %q want %q", got.Title(), post.Title())
	}

	if got.Content() != post.Content() {
		t.Errorf("Content mismatch: got %q want %q", got.Content(), post.Content())
	}

	if got.AuthorID() != author.ID() {
		t.Errorf("AuthorID mismatch: got %v want %v", got.AuthorID(), author.ID())
	}

	if len(got.Tags()) != len(tags) {
		t.Fatalf("Tags length mismatch: got %d want %d", len(got.Tags()), len(tags))
	}

	if len(got.Stances()) != 3 {
		js, _ := json.Marshal(got.Stances())
		t.Fatalf("expected 3 stances, got %d: %s", len(got.Stances()), js)
	}
}

func TestPostRepositoryPG_ListByUserID(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	repo := postgres.NewPostRepositoryPG(testPool)

	now := time.Now().UTC()

	author, _ := domain.NewUser(now, "listpostsuser", "Lister")
	userRepo := postgres.NewUserRepositoryPG(testPool)
	if err := userRepo.Create(ctx, author, "hash"); err != nil {
		t.Fatalf("user create: %v", err)
	}

	color1, _ := domain.NewColor("#123456")

	makePost := func(title string) *domain.Post {
		s1, _ := domain.NewStance(uuid.New(), "S1", color1, "")
		s2, _ := domain.NewStance(uuid.New(), "S2", color1, "")
		s3, _ := domain.NewStance(uuid.New(), "S3", color1, "")
		p, err := domain.NewPost(
			now,
			title,
			"Body",
			nil,
			[]domain.Stance{s1, s2, s3},
			author.ID(),
		)
		if err != nil {
			t.Fatalf("NewPost error: %v", err)
		}
		return p
	}

	post1 := makePost("Hello")
	post2 := makePost("World")

	if err := repo.Create(ctx, post1); err != nil {
		t.Fatalf("create p1: %v", err)
	}
	if err := repo.Create(ctx, post2); err != nil {
		t.Fatalf("create p2: %v, %s, %s", err, post1.ID().String(), post2.ID().String())
	}

	page := filters.NewPage(10, 0)
	list, err := repo.ListByUserID(ctx, author.ID(), page)
	if err != nil {
		t.Fatalf("ListByUserID error: %v", err)
	}

	if len(list) < 2 {
		t.Fatalf("expected at least 2 posts, got %d", len(list))
	}

	titles := []string{list[0].Title(), list[1].Title()}
	if !(titles[0] == "Hello" || titles[1] == "World") {
		t.Errorf("unexpected titles: %v", titles)
	}
}

func TestPostRepositoryPG_Delete(t *testing.T) {
	if testPool == nil {
		t.Skip("No test database configured")
	}

	ctx := context.Background()
	repo := postgres.NewPostRepositoryPG(testPool)

	now := time.Now().UTC()

	author, _ := domain.NewUser(now, "delpost", "DelUser")
	userRepo := postgres.NewUserRepositoryPG(testPool)
	if err := userRepo.Create(ctx, author, "hash"); err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

	color, _ := domain.NewColor("#abcdef")
	s1, _ := domain.NewStance(uuid.New(), "S1", color, "")
	s2, _ := domain.NewStance(uuid.New(), "S2", color, "")
	s3, _ := domain.NewStance(uuid.New(), "S2", color, "")
	post, err := domain.NewPost(now, "To delete", "body", nil, []domain.Stance{s1, s2, s3}, author.ID())
	if err != nil {
		t.Fatalf("NewPost error: %v", err)
	}

	if err := repo.Create(ctx, post); err != nil {
		t.Fatalf("create post: %v", err)
	}

	if err := repo.Delete(ctx, post.ID()); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = repo.GetByID(ctx, post.ID())
	if err == nil {
		t.Fatalf("expected error for deleted post, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
