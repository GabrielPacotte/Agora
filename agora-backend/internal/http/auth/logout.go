package authhttp

import (
	"net/http"
	"time"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := r.Cookie("refresh_token")
	if err == nil && c.Value != "" {
		_ = h.service.Logout(ctx, c.Value)
	}

	clearRefreshCookie(w, h.secureCookie, h.sameSite)
	w.WriteHeader(http.StatusNoContent)
}

func clearRefreshCookie(w http.ResponseWriter, secure bool, sameSite httpSameSite) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/v1/auth",
		HttpOnly: true,
		Secure:   secure,
		MaxAge:   0,
		Expires:  time.Unix(0, 0),
	}

	switch sameSite {
	case sameSiteLax:
		c.SameSite = http.SameSiteLaxMode
	case sameSiteStrict:
		c.SameSite = http.SameSiteStrictMode
	case sameSiteNone:
		c.SameSite = http.SameSiteNoneMode
	default:
		c.SameSite = http.SameSiteDefaultMode
	}

	http.SetCookie(w, c)
}
