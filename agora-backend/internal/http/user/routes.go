package userhttp

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/me", h.GetMe)
	r.Delete("/me", h.DeleteMe)

	r.Get("/me/preferences", h.GetPreferences)
	r.Put("/me/preferences", h.SetPreferences)

	r.Get("/me/saved", h.ListSavedContent)
	r.Post("/me/saved", h.SaveContent)
	r.Delete("/me/saved", h.UnsaveContent)

	return r
}