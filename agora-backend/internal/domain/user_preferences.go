package domain

import (
	"github.com/google/uuid"
)

type UserPreferences struct {
	userID uuid.UUID
	tags   []Tag
}

func NewUserPreferences(userID uuid.UUID, tags []Tag) (*UserPreferences, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidID
	}
	tags = dedupeTags(tags)
	return &UserPreferences{
		userID: userID,
		tags:   tags,
	}, nil
}

func (p *UserPreferences) UserID() uuid.UUID {
	return p.userID
}

func (p *UserPreferences) Tags() []Tag {
	out := make([]Tag, len(p.tags))
	copy(out, p.tags)
	return out
}

func (p *UserPreferences) SetTags(tags []Tag) {
	p.tags = dedupeTags(tags)
}
