package domain_test

import (
	"strings"
	"testing"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

var NewTag = domain.NewTag
var ErrInvalidTag = domain.ErrInvalidTag

func TestNewTag_ValidTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  domain.Tag
	}{
		{"simple word", "politics", domain.Tag("politics")},
		{"numbers", "123", domain.Tag("123")},
		{"with dash", "video-games", domain.Tag("video-games")},
		{"with underscore", "game_dev", domain.Tag("game_dev")},
		{"leading/trailing spaces trimmed", "   science   ", domain.Tag("science")},
		{"uppercase converted to lowercase", "FrAnCe", domain.Tag("france")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTag(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestNewTag_InvalidTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"spaces only", "   "},
		{"invalid characters", "hello!"},
		{"starts with dash", "-tag"},
		{"ends with underscore", "tag_"},
		{"too long", "a" + strings.Repeat("b", 40)},
		{"spaces inside", "my tag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTag(tt.input)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if err != ErrInvalidTag {
				t.Fatalf("expected ErrInvalidTag, got %v", err)
			}
		})
	}
}

func TestTag_String(t *testing.T) {
	tag := domain.Tag("example")
	if tag.String() != "example" {
		t.Fatalf("expected %q, got %q", "example", tag.String())
	}
}
