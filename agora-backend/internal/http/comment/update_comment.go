package commenthttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "commentId")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid comment id")
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	requesterID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	updated, err := h.commentService.Update(ctx, requesterID, commentID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "comment not found")
		case errors.Is(err, domain.ErrForbidden):
			httpcommon.WriteError(w, http.StatusForbidden, httpcommon.ErrCodeForbidden, "forbidden")
		case errors.Is(err, domain.ErrEmptyCommentContent),
			errors.Is(err, domain.ErrCommentContentTooLong):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, err.Error())
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	resp := UpdateCommentResponse{
		Comment: mapComment(updated),
		Meta:    EmptyMeta{},
	}
	httpcommon.WriteJSON(w, http.StatusOK, resp)
}

func mapComment(c *domain.Comment) CommentJSON {
	var replyTo *string
	if c.ReplyToID() != nil {
		s := c.ReplyToID().String()
		replyTo = &s
	}

	return CommentJSON{
		ID:        c.ID().String(),
		PostID:    c.PostID().String(),
		AuthorID:  c.AuthorID().String(),
		ReplyToID: replyTo,
		Content:   c.Content(),
		StanceID:  c.Stance().ID().String(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}
