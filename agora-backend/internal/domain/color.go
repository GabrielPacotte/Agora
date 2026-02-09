package domain

import (
	"regexp"
	"strings"
)

type Color string

var colorRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func NewColor(hexaCode string) (Color, error) {
	hexaCode = strings.TrimSpace(hexaCode)
	hexaCode = strings.ToLower(hexaCode)
	if !colorRegex.MatchString(hexaCode) {
		return "", ErrInvalidColor
	}
	return Color(hexaCode), nil
}

func (c Color) String() string {
	return string(c)
}
