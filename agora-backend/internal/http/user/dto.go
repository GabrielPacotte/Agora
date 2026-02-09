package userhttp

import (
	"time"
)

type EmptyMeta struct{}

type UserJSON struct {
	ID          string    `json:"id"`
	Login       string    `json:"login"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetMeResponse struct {
	User UserJSON  `json:"user"`
	Meta EmptyMeta `json:"meta"`
}

type GetPreferencesResponse struct {
	Tags []string  `json:"tags"`
	Meta EmptyMeta `json:"meta"`
}

type SetPreferencesRequest struct {
	Tags []string `json:"tags"`
}

type SetPreferencesResponse struct {
	Meta EmptyMeta `json:"meta"`
}

type SavedContentJSON struct {
	UserID      string    `json:"user_id"`
	ContentID   string    `json:"content_id"`
	ContentType string    `json:"content_type"`
	SavedAt     time.Time `json:"saved_at"`
}

type PageMeta struct {
	NextOffset *int `json:"next_offset,omitempty"`
	Count      int  `json:"count"`
}

type ListSavedContentResponse struct {
	Saved []SavedContentJSON `json:"saved"`
	Meta  PageMeta           `json:"meta"`
}

type SaveContentRequest struct {
	ContentID   string `json:"content_id"`
	ContentType string `json:"content_type"`
}

type SaveContentResponse struct {
	Meta EmptyMeta `json:"meta"`
}

type UnsaveContentResponse struct {
	Meta EmptyMeta `json:"meta"`
}
