package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func NewCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://agora.example",
			"http://localhost:3000",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		ExposedHeaders: []string{
			"X-Request-ID",
		},

		AllowCredentials: true,

		MaxAge: 300,
	})
}
