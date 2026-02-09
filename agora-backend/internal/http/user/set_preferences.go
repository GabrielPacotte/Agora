package userhttp

import (
	"encoding/json"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) SetPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	var req SetPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	tags := make([]domain.Tag, 0, len(req.Tags))
	for _, raw := range req.Tags {
		t, err := domain.NewTag(raw)
		if err != nil {
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid tag")
			return
		}
		tags = append(tags, t)
	}

	if err := h.userService.SetPreferences(ctx, userID, tags); err != nil {
		httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, SetPreferencesResponse{Meta: EmptyMeta{}})
}
