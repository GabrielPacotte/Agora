package authhttp

import (
	"github.com/GabrielPacotte/Agora/internal/services"
)

type Handler struct {
	service      services.AuthService
	secureCookie bool
	sameSite     httpSameSite
}

type httpSameSite int

const (
	sameSiteDefault httpSameSite = iota
	sameSiteLax
	sameSiteStrict
	sameSiteNone
)

func NewHandler(
	service services.AuthService,
	secureCookie bool,
) *Handler {
	return &Handler{
		service:      service,
		secureCookie: secureCookie,
		sameSite:     sameSiteLax,
	}
}
