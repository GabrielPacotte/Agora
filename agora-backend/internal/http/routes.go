package httpserver

import (
	"net/http"

	authhttp "github.com/GabrielPacotte/Agora/internal/http/auth"
	commenthttp "github.com/GabrielPacotte/Agora/internal/http/comment"
	custommiddleware "github.com/GabrielPacotte/Agora/internal/http/middleware"
	posthttp "github.com/GabrielPacotte/Agora/internal/http/post"
	userhttp "github.com/GabrielPacotte/Agora/internal/http/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	Auth     *authhttp.Handler
	Posts    *posthttp.Handler
	Comments *commenthttp.Handler
	Users    *userhttp.Handler

	AuthMw func(h http.Handler) http.Handler
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.NewCORS())

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", authhttp.Routes(h.Auth))

		r.Group(func(r chi.Router) {
			r.Use(h.AuthMw)

			r.Mount("/posts", posthttp.Routes(h.Posts))
			r.Mount("/comments", commenthttp.Routes(h.Comments))
			r.Mount("/users", userhttp.Routes(h.Users))
		})
	})

	return r
}
