package userhttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/google/uuid"
)

func (h *Handler) UnsaveContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	contentIDStr := r.URL.Query().Get("content_id")
	contentTypeStr := r.URL.Query().Get("content_type")

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid content_id")
		return
	}

	ct := domain.ContentType(contentTypeStr)
	if ct != domain.POST && ct != domain.COMMENT {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid content_type")
		return
	}

	err = h.userService.UnsaveContent(ctx, userID, contentID, ct)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "saved content not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
