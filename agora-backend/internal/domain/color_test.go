package domain_test

import (
	"testing"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

func TestNewColor_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  domain.Color
	}{
		{"lowercase hex", "#abcdef", domain.Color("#abcdef")},
		{"uppercase hex", "#ABCDEF", domain.Color("#abcdef")},
		{"mixed case", "#AbCdEf", domain.Color("#abcdef")},
		{"trimmed", "   #123456   ", domain.Color("#123456")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewColor(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestNewColor_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"no # prefix", "123456"},
		{"too short", "#abc"},
		{"too long", "#1234567"},
		{"invalid char", "#12gg56"},
		{"spaces only", "    "},
		{"missing hex", "#------"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewColor(tt.input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidColor {
				t.Fatalf("expected ErrInvalidColor, got %v", err)
			}
		})
	}
}

func TestColor_String(t *testing.T) {
	c := domain.Color("#aabbcc")
	if c.String() != "#aabbcc" {
		t.Fatalf("expected %q, got %q", "#aabbcc", c.String())
	}
}
