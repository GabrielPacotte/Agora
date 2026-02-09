package domain

import (
	"time"

	"github.com/google/uuid"
)

type SavedContent struct {
	userID      uuid.UUID
	subjectID   uuid.UUID
	subjectType ContentType
	savedAt     time.Time
}

func NewSavedContent(savedAt time.Time, userID, subjectID uuid.UUID, subjectType ContentType) (*SavedContent, error) {
	if userID == uuid.Nil || subjectID == uuid.Nil {
		return nil, ErrInvalidID
	}

	if subjectType != POST && subjectType != COMMENT {
		return nil, ErrInvalidContentType
	}

	return &SavedContent{
		userID:      userID,
		subjectID:   subjectID,
		subjectType: subjectType,
		savedAt:     savedAt,
	}, nil
}

func (s *SavedContent) UserID() uuid.UUID {
	return s.userID
}

func (s *SavedContent) SubjectID() uuid.UUID {
	return s.subjectID
}

func (s *SavedContent) SubjectType() ContentType {
	return s.subjectType
}

func (s *SavedContent) SavedAt() time.Time {
	return s.savedAt
}
