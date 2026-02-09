package posthttp

import (
	"time"
)

type CreateStanceRequest struct {
	Label       string `json:"label"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type CreatePostRequest struct {
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Tags    []string              `json:"tags"`
	Stances []CreateStanceRequest `json:"stances"`
}

type StanceJSON struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type PostJSON struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Tags      []string     `json:"tags"`
	Stances   []StanceJSON `json:"stances"`
	AuthorID  string       `json:"author_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type EmptyMeta struct {
}

type GetPostResponse struct {
	Post PostJSON  `json:"post"`
	Meta EmptyMeta `json:"meta"`
}

type UpdateStanceRequest struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type UpdatePostRequest struct {
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Tags    []string              `json:"tags"`
	Stances []UpdateStanceRequest `json:"stances"`
}

type CreateCommentRequest struct {
	StanceID  string  `json:"stance_id"`
	Content   string  `json:"content"`
	ReplyToID *string `json:"reply_to_id"`
}

type CommentJSON struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	ReplyToID *string   `json:"reply_to_id,omitempty"`
	Content   string    `json:"content"`
	StanceID  string    `json:"stance_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentResponse struct {
	Comment CommentJSON `json:"comment"`
	Meta    EmptyMeta   `json:"meta"`
}

type PageMeta struct {
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	NextOffset *int `json:"next_offset,omitempty"`
	Count      int  `json:"count"`
}

type ListRootCommentsResponse struct {
	Comments []CommentJSON `json:"comments"`
	Meta     PageMeta      `json:"meta"`
}

type ListPostsResponse struct {
	Posts []PostJSON `json:"posts"`
	Meta  PageMeta   `json:"meta"`
}
