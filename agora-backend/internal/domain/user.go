package domain

import (
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type User struct {
	identified
	login       string
	displayName string
}

var loginRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9_-]{1,30})[a-z0-9]$`)

func NewUser(date time.Time, login, displayName string) (*User, error) {
	identified, err := newIdentified(date)
	if err != nil {
		return nil, err
	}
	login = strings.ToLower(strings.TrimSpace(login))
	displayName = strings.TrimSpace(displayName)
	err = validateUserFields(login, displayName)
	if err != nil {
		return nil, err
	}
	return &User{
		identified:  identified,
		login:       login,
		displayName: displayName,
	}, nil
}

func ExistingUser(id uuid.UUID, createdAt, updatedAt time.Time, login, displayName string) (*User, error) {
	identified, err := existingIdentified(id, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	login = strings.ToLower(strings.TrimSpace(login))
	displayName = strings.TrimSpace(displayName)
	err = validateUserFields(login, displayName)
	if err != nil {
		return nil, err
	}
	return &User{
		identified:  identified,
		login:       login,
		displayName: displayName,
	}, nil
}

func validateUserFields(login, displayName string) error {
	if !isValidLogin(login) {
		return ErrInvalidLogin
	}
	if !isValidDisplayName(displayName) {
		return ErrInvalidDisplayName
	}
	return nil
}

func isValidLogin(login string) bool {
	return loginRegex.MatchString(login)
}

func isValidDisplayName(displayName string) bool {
	if len(displayName) == 0 || utf8.RuneCountInString(displayName) > 50 {
		return false
	}

	if strings.TrimSpace(displayName) == "" {
		return false
	}

	for _, r := range displayName {
		if unicode.IsControl(r) {
			return false
		}

		switch {
		case unicode.IsLetter(r),
			unicode.IsDigit(r),
			unicode.IsSpace(r),
			r == '-',
			r == '\'',
			r == '.':
		default:
			return false
		}
	}

	return true
}

func (u *User) Login() string {
	return u.login
}

func (u *User) DisplayName() string {
	return u.displayName
}
