package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxCommentLength = 5000

type Comment struct {
	identified
	content   string
	stance    Stance
	postID    uuid.UUID
	authorID  uuid.UUID
	replyToID *uuid.UUID
}

func NewComment(date time.Time, content string, stance Stance, postID, authorID uuid.UUID, replyToID *uuid.UUID) (*Comment, error) {
	identified, err := newIdentified(date)
	if err != nil {
		return nil, err
	}
	err = validateCommentContent(content)
	if err != nil {
		return nil, err
	}
	if authorID == uuid.Nil || postID == uuid.Nil || (replyToID != nil && replyToID == &uuid.Nil) {
		return nil, ErrInvalidID
	}
	return &Comment{
		identified: identified,
		content:    content,
		stance:     stance,
		postID:     postID,
		replyToID:  replyToID,
		authorID:   authorID,
	}, nil
}

func ExistingComment(id uuid.UUID, createdAt, updatedAt time.Time, content string, stance Stance, postID, authorID uuid.UUID, replyToID *uuid.UUID) (*Comment, error) {
	identified, err := existingIdentified(id, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	err = validateCommentContent(content)
	if err != nil {
		return nil, err
	}
	if authorID == uuid.Nil || postID == uuid.Nil || (replyToID != nil && *replyToID == uuid.Nil) {
		return nil, ErrInvalidID
	}
	return &Comment{
		identified: identified,
		content:    content,
		stance:     stance,
		postID:     postID,
		replyToID:  replyToID,
		authorID:   authorID,
	}, nil
}

func validateCommentContent(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return ErrEmptyCommentContent
	}
	if utf8.RuneCountInString(content) > maxCommentLength {
		return ErrCommentContentTooLong
	}
	return nil
}

func (c *Comment) Content() string {
	return c.content
}

func (c *Comment) Stance() Stance {
	return c.stance
}

func (c *Comment) PostID() uuid.UUID {
	return c.postID
}

func (c *Comment) ReplyToID() *uuid.UUID {
	return c.replyToID
}

func (c *Comment) IsReply() bool {
	return c.replyToID != nil
}

func (c *Comment) AuthorID() uuid.UUID {
	return c.authorID
}
