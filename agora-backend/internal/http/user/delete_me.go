package userhttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	err := h.userService.DeleteAccount(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "user not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
