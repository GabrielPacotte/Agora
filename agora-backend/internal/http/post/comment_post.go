package posthttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GabrielPacotte/Agora/internal/domain"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postIDStr := chi.URLParam(r, "postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid post id")
		return
	}

	requesterID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid json body")
		return
	}

	stanceID, err := uuid.Parse(req.StanceID)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid stance_id")
		return
	}

	var replyToID *uuid.UUID
	if req.ReplyToID != nil {
		id, err := uuid.Parse(*req.ReplyToID)
		if err != nil {
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid reply_to_id")
			return
		}
		replyToID = &id
	}

	comment, err := h.commentService.Create(ctx, requesterID, postID, stanceID, req.Content, replyToID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "post, stance or comment to reply to not found")
		case errors.Is(err, domain.ErrForbidden):
			httpcommon.WriteError(w, http.StatusForbidden, httpcommon.ErrCodeForbidden, "forbidden")
		case errors.Is(err, domain.ErrEmptyCommentContent),
			errors.Is(err, domain.ErrCommentContentTooLong),
			errors.Is(err, domain.ErrInvalidID):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, err.Error())
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	resp := CreateCommentResponse{
		Comment: mapComment(comment),
		Meta:    EmptyMeta{},
	}

	httpcommon.WriteJSON(w, http.StatusCreated, resp)
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
