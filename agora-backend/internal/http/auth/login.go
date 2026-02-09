package authhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	pair, _, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials),
			errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "invalid credentials")
		case errors.Is(err, domain.ErrInvalidPassword),
			errors.Is(err, domain.ErrInvalidLogin):
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

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
