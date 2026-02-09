package domain

import "errors"

var (
	// Model errors

	ErrInvalidTimestamp                = errors.New("invalid timestamp")
	ErrInvalidID                       = errors.New("invalid UUID")
	ErrInvalidLogin                    = errors.New("invalid user login")
	ErrInvalidDisplayName              = errors.New("invalid user display name")
	ErrInvalidColor                    = errors.New("invalid color hexadecimal code")
	ErrInvalidTag                      = errors.New("invalid tag value")
	ErrMinimalStanceAmount             = errors.New("a Post must have at least 3 stances")
	ErrEmptyReasons                    = errors.New("reasons cannot be empty")
	ErrEmptyStanceLabel                = errors.New("stance label cannot be an empty string")
	ErrEmptyCommentContent             = errors.New("comment content cannot be empty")
	ErrCommentContentTooLong           = errors.New("comment content is too long")
	ErrEmptyPostTitle                  = errors.New("post title cannot be empty")
	ErrPostTitleTooLong                = errors.New("post title is too long")
	ErrPostContentTooLong              = errors.New("post description is too long")
	ErrReportDescriptionTooLong        = errors.New("report description is too long")
	ErrInvalidContentType              = errors.New("invalid content type value for report")
	ErrInvalidModerationActionDecision = errors.New("invalid moderation action decision value")
	ErrModerationActionDescriptionSize = errors.New("invalid size for moderation action decision description")

	// Repositories errors

	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrValidation          = errors.New("validation error")
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrTimeout             = errors.New("timeout")
	ErrCanceled            = errors.New("canceled")
	ErrForeignKeyViolation = errors.New("foreign key violation")

	// Auth / credentials
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("invalid password")

	// Tokens / sessions
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)
