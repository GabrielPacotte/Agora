package userhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/google/uuid"
)

func (h *Handler) SaveContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	var req SaveContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	contentID, err := uuid.Parse(req.ContentID)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid content_id")
		return
	}

	ct := domain.ContentType(req.ContentType)
	if ct != domain.POST && ct != domain.COMMENT {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid content_type")
		return
	}

	if err := h.userService.SaveContent(ctx, userID, contentID, ct); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "content not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, SaveContentResponse{Meta: EmptyMeta{}})
}
