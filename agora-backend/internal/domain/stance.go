package domain

import (
	"strings"

	"github.com/google/uuid"
)

type Stance struct {
	id          uuid.UUID
	label       string
	color       Color
	description string
}

func NewStance(id uuid.UUID, label string, color Color, desc string) (Stance, error) {
	label = strings.TrimSpace(label)
	desc = strings.TrimSpace(desc)
	if label == "" {
		return Stance{}, ErrEmptyStanceLabel
	}
	return Stance{
		id:          id,
		label:       label,
		color:       color,
		description: desc,
	}, nil
}

func (s Stance) ID() uuid.UUID {
	return s.id
}

func (s Stance) Label() string {
	return s.label
}

func (s Stance) Color() Color {
	return s.color
}

func (s Stance) Description() string {
	return s.description
}

func (s Stance) HasDescription() bool {
	return s.description != ""
}
