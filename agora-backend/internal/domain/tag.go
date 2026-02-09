package domain

import (
	"regexp"
	"strings"
)

type Tag string

var tagRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9_-]{0,30}[a-z0-9])?$`)

func NewTag(value string) (Tag, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	if !tagRegex.MatchString(value) {
		return "", ErrInvalidTag
	}
	return Tag(value), nil
}

func (t Tag) String() string {
	return string(t)
}
