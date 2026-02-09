package commenthttp

import "time"

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

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type UpdateCommentResponse struct {
	Comment CommentJSON `json:"comment"`
	Meta    EmptyMeta   `json:"meta"`
}

type ListRepliesResponse struct {
	Comments []CommentJSON `json:"comments"`
	Meta     PageMeta      `json:"meta"`
}

type PageMeta struct {
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	NextOffset *int `json:"next_offset,omitempty"`
	Count      int  `json:"count"`
}

type EmptyMeta struct {
}
