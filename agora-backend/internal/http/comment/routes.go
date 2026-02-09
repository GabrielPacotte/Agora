package commenthttp

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{commentId}/replies", h.ListReplies)
	r.Delete("/{commentId}", h.DeleteComment)
	r.Patch("/{commentId}", h.UpdateComment)
	return r
}
