package userhttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	u, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "user not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	resp := GetMeResponse{
		User: UserJSON{
			ID:          u.ID().String(),
			Login:       u.Login(),
			DisplayName: u.DisplayName(),
			CreatedAt:   u.CreatedAt(),
			UpdatedAt:   u.UpdatedAt(),
		},
		Meta: EmptyMeta{},
	}

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
