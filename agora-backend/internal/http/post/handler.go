package posthttp

import (
	"github.com/GabrielPacotte/Agora/internal/services"
)

type Handler struct {
	postService    services.PostService
	commentService services.CommentService
}

func NewHandler(postService services.PostService, commentService services.CommentService) *Handler {
	return &Handler{
		postService:    postService,
		commentService: commentService,
	}
}
