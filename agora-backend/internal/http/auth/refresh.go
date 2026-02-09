package authhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "missing refresh token")
		return
	}

	pair, _, err := h.service.Refresh(ctx, c.Value)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken),
			errors.Is(err, domain.ErrUnauthorized),
			errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "invalid refresh token")
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

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
