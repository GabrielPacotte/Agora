package posthttp

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) ListRootCommentsForPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postIDStr := chi.URLParam(r, "postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid post id")
		return
	}

	limit, offset, err := parseLimitOffset(r, 20, 0)
	if err != nil {
		httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.ErrCodeInvalidInput, "invalid pagination")
		return
	}

	page := filters.NewPage(limit, offset)

	comments, err := h.commentService.ListRootComments(ctx, postID, page)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.ErrCodeNotFound, "post not found")
		default:
			httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		}
		return
	}

	out := make([]CommentJSON, 0, len(comments))
	for _, c := range comments {
		out = append(out, mapComment(&c))
	}

	meta := PageMeta{
		Limit:  limit,
		Offset: offset,
		Count:  len(out),
	}

	if len(out) == limit {
		next := offset + limit
		meta.NextOffset = &next
	}

	resp := ListRootCommentsResponse{
		Comments: out,
		Meta:     meta,
	}

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}

func parseLimitOffset(r *http.Request, defaultLimit, defaultOffset int) (int, int, error) {
	limit := defaultLimit
	offset := defaultOffset

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, err
		}
		limit = n
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, err
		}
		offset = n
	}

	if limit < 1 || limit > 100 {
		return 0, 0, errors.New("limit out of range")
	}
	if offset < 0 {
		return 0, 0, errors.New("offset out of range")
	}

	return limit, offset, nil
}
