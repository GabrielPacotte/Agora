package commenthttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "commentId")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid comment id")
		return
	}

	requesterID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	err = h.commentService.Delete(ctx, requesterID, commentID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "comment not found")
		case errors.Is(err, domain.ErrForbidden):
			httpcommon.WriteError(w, http.StatusForbidden, httpcommon.ErrCodeForbidden, "forbidden")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
