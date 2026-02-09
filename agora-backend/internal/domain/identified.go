package domain

import (
	"time"

	"github.com/google/uuid"
)

type identified struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func newIdentified(date time.Time) (identified, error) {
	if date.IsZero() {
		return identified{}, ErrInvalidTimestamp
	}
	return identified{
		id:        uuid.New(),
		createdAt: date,
		updatedAt: date,
	}, nil
}

func existingIdentified(id uuid.UUID, createdAt, updatedAt time.Time) (identified, error) {
	if id == uuid.Nil {
		return identified{}, ErrInvalidID
	}
	if createdAt.IsZero() || updatedAt.IsZero() || updatedAt.Before(createdAt) {
		return identified{}, ErrInvalidTimestamp
	}
	return identified{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (i *identified) ID() uuid.UUID {
	return i.id
}

func (i *identified) CreatedAt() time.Time {
	return i.createdAt
}

func (i *identified) UpdatedAt() time.Time {
	return i.updatedAt
}
