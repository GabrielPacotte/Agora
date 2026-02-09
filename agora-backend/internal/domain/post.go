package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const MaxPostTitleLength = 200
const MaxPostContentLength = 10000

type Post struct {
	identified
	title    string
	content  string
	tags     []Tag
	stances  []Stance
	authorID uuid.UUID
}

func NewPost(date time.Time, title string, content string, tags []Tag, stances []Stance, authorID uuid.UUID) (*Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	tags = dedupeTags(tags)
	stances = dedupeStances(stances)
	identified, err := newIdentified(date)
	if err != nil {
		return nil, err
	}
	err = validatePostFields(title, content, stances, authorID)
	if err != nil {
		return nil, err
	}
	return &Post{
		identified: identified,
		title:      title,
		content:    content,
		tags:       tags,
		stances:    stances,
		authorID:   authorID,
	}, nil
}

func ExistingPost(id uuid.UUID, createdAt, updatedAt time.Time, title, content string, tags []Tag, stances []Stance, authorID uuid.UUID) (*Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	tags = dedupeTags(tags)
	stances = dedupeStances(stances)
	identified, err := existingIdentified(id, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	err = validatePostFields(title, content, stances, authorID)
	if err != nil {
		return nil, err
	}
	return &Post{
		identified: identified,
		title:      title,
		content:    content,
		tags:       tags,
		stances:    stances,
		authorID:   authorID,
	}, nil
}

func dedupeTags(tags []Tag) []Tag {
	seen := make(map[Tag]struct{}, len(tags))
	out := make([]Tag, 0, len(tags))

	for _, t := range tags {
		if _, exists := seen[t]; exists {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func dedupeStances(stances []Stance) []Stance {
	seen := make(map[uuid.UUID]struct{}, len(stances))
	out := make([]Stance, 0, len(stances))

	for _, s := range stances {
		if _, exists := seen[s.ID()]; exists {
			continue
		}
		seen[s.ID()] = struct{}{}
		out = append(out, s)
	}

	return out
}

func validatePostFields(title, content string, stances []Stance, authorID uuid.UUID) error {
	if title == "" {
		return ErrEmptyPostTitle
	}
	if utf8.RuneCountInString(title) > MaxPostTitleLength {
		return ErrPostTitleTooLong
	}

	if utf8.RuneCountInString(content) > MaxPostContentLength {
		return ErrPostContentTooLong
	}

	if len(stances) < 3 {
		return ErrMinimalStanceAmount
	}

	if authorID == uuid.Nil {
		return ErrInvalidID
	}

	return nil
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) HasContent() bool {
	return p.content != ""
}

func (p *Post) Tags() []Tag {
	out := make([]Tag, len(p.tags))
	copy(out, p.tags)
	return out
}

func (p *Post) Stances() []Stance {
	out := make([]Stance, len(p.stances))
	copy(out, p.stances)
	return out
}

func (p *Post) AuthorID() uuid.UUID {
	return p.authorID
}
