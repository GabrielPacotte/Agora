package posthttp

import (
	"github.com/go-chi/chi/v5"
)

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreatePost)
	r.Get("/", h.GetFeed)
	r.Get("/search", h.Search)

	r.Get("/{postId}", h.GetPostByID)
	r.Delete("/{postId}", h.DeletePost)
	r.Patch("/{postId}", h.UpdatePost)

	r.Post("/{postId}/comments", h.CreateComment)
	r.Get("/{postId}/comments", h.ListRootCommentsForPost)
	return r
}
