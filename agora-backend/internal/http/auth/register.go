package authhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	pair, _, err := h.service.Register(ctx, req.Login, req.DisplayName, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAlreadyExists):
			httpcommon.WriteError(w, http.StatusConflict, httpcommon.ErrCodeConflict, "login already exists")
		case errors.Is(err, domain.ErrInvalidPassword),
			errors.Is(err, domain.ErrInvalidLogin),
			errors.Is(err, domain.ErrInvalidDisplayName):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, err.Error())
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	expiresAt := pair.Refresh.ExpiresAt()
	setRefreshCookie(w, pair.Refresh.Raw(), time.Until(expiresAt), h.secureCookie, h.sameSite)

	resp := AuthResponse{
		Auth: AuthJSON{
			AccessToken: string(pair.Access),
			TokenType:   "Bearer",
			ExpiresIn:   int64(pair.ExpiresIn.Seconds()),
		},
		Meta: EmptyMeta{},
	}

	httpcommon.WriteJSON(w, http.StatusCreated, resp)
}

func setRefreshCookie(w http.ResponseWriter, raw string, ttl time.Duration, secure bool, sameSite httpSameSite) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    raw,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		MaxAge:   int(ttl.Seconds()),
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
