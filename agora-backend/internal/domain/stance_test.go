package domain_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/GabrielPacotte/Agora/internal/domain"
)

func TestNewStance_Valid(t *testing.T) {
	id := uuid.New()
	color, _ := domain.NewColor("#abcdef")

	tests := []struct {
		name  string
		label string
		desc  string
		color domain.Color
	}{
		{"simple label", "Agree", "", color},
		{"trimmed label", "   Neutral   ", "desc", color},
		{"with description", "Oppose", "Explanation", color},
		{"trimmed description", "Viewpoint", "   Something   ", color},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := domain.NewStance(id, tt.label, tt.color, tt.desc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check ID
			if s.ID() != id {
				t.Fatalf("ID mismatch: expected %v, got %v", id, s.ID())
			}

			// Check trimmed label
			if s.Label() != strings.TrimSpace(tt.label) {
				t.Fatalf("Label mismatch: expected %q, got %q",
					strings.TrimSpace(tt.label), s.Label())
			}

			// Check Color
			if s.Color() != tt.color {
				t.Fatalf("Color mismatch: expected %v, got %v", tt.color, s.Color())
			}

			// Check trimmed description
			if s.Description() != strings.TrimSpace(tt.desc) {
				t.Fatalf("Description mismatch: expected %q, got %q",
					strings.TrimSpace(tt.desc), s.Description())
			}
		})
	}
}

func TestNewStance_Invalid(t *testing.T) {
	color, _ := domain.NewColor("#123456")
	id := uuid.New()

	tests := []struct {
		name  string
		label string
	}{
		{"empty label", ""},
		{"spaces only", "     "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewStance(id, tt.label, color, "desc")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrEmptyStanceLabel {
				t.Fatalf("expected ErrEmptyStanceLabel, got %v", err)
			}
		})
	}
}

func TestStance_HasDescription(t *testing.T) {
	color, _ := domain.NewColor("#aaaaaa")
	id := uuid.New()

	s1, _ := domain.NewStance(id, "Agree", color, "")
	if s1.HasDescription() != false {
		t.Fatalf("expected HasDescription=false for empty description")
	}

	s2, _ := domain.NewStance(id, "Agree", color, "Hello")
	if s2.HasDescription() != true {
		t.Fatalf("expected HasDescription=true for non-empty description")
	}
}
