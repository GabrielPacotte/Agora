package posthttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "postId")
	postID, err := uuid.Parse(idStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid id")
		return
	}

	requesterID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	if err := h.postService.Delete(ctx, requesterID, postID); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "post not found")
		case errors.Is(err, domain.ErrForbidden):
			httpcommon.WriteError(w, http.StatusForbidden, httpcommon.ErrCodeForbidden, "forbidden")
		case errors.Is(err, domain.ErrInvalidID):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid input")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
