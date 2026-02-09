package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type ContentType string

const (
	POST    ContentType = "post"
	COMMENT ContentType = "comment"
)

const maxReportDescriptionLength = 2000

type Report struct {
	identified
	subjectID   uuid.UUID
	subjectType ContentType
	authorID    uuid.UUID
	description string
}

func NewReport(date time.Time, authorID, subjectID uuid.UUID, subjectType ContentType, desc string) (*Report, error) {
	identified, err := newIdentified(date)
	if err != nil {
		return nil, err
	}
	desc = strings.TrimSpace(desc)
	err = validateReportFields(authorID, subjectID, subjectType, desc)
	if err != nil {
		return nil, err
	}
	return &Report{
		identified:  identified,
		subjectID:   subjectID,
		subjectType: subjectType,
		authorID:    authorID,
		description: desc,
	}, nil
}

func ExistingReport(id uuid.UUID, createdAt, updatedAt time.Time, authorID, subjectID uuid.UUID, subjectType ContentType, desc string) (*Report, error) {
	ident, err := existingIdentified(id, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	desc = strings.TrimSpace(desc)
	if err := validateReportFields(authorID, subjectID, subjectType, desc); err != nil {
		return nil, err
	}

	return &Report{
		identified:  ident,
		subjectID:   subjectID,
		subjectType: subjectType,
		authorID:    authorID,
		description: desc,
	}, nil
}

func validateReportFields(authorID, subjectID uuid.UUID, subjectType ContentType, desc string) error {
	if subjectID == uuid.Nil || authorID == uuid.Nil {
		return ErrInvalidID
	}

	if subjectType != POST && subjectType != COMMENT {
		return ErrInvalidContentType
	}

	if utf8.RuneCountInString(desc) > maxReportDescriptionLength {
		return ErrReportDescriptionTooLong
	}
	return nil
}

func (r *Report) SubjectID() uuid.UUID {
	return r.subjectID
}

func (r *Report) SubjectType() ContentType {
	return r.subjectType
}

func (r *Report) AuthorID() uuid.UUID {
	return r.authorID
}

func (r *Report) Description() string {
	return r.description
}

func (r *Report) HasDescription() bool {
	return r.description != ""
}
