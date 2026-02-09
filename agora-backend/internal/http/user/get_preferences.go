package userhttp

import (
	"net/http"

	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	tags, err := h.userService.ListPreferences(ctx, userID)
	if err != nil {
		httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		return
	}

	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.String())
	}

	resp := GetPreferencesResponse{
		Tags: out,
		Meta: EmptyMeta{},
	}
	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
