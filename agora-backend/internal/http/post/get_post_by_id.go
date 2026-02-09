package posthttp

import (
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "postId")
	postID, err := uuid.Parse(idStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid id")
		return
	}

	post, err := h.postService.GetPostByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "post not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	resp := mapPostToResponse(post)

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
