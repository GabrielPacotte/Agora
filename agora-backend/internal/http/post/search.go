package posthttp

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	var tags []domain.Tag
	if raw := q.Get("tags"); raw != "" {
		for _, t := range strings.Split(raw, ",") {
			tags = append(tags, domain.Tag(strings.TrimSpace(t)))
		}
	}

	var fromDate *time.Time
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			fromDate = &t
		}
	}

	var toDate *time.Time
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			toDate = &t
		}
	}

	filter := filters.SearchPostFilter{
		Page:     filters.NewPage(limit, offset),
		Tags:     tags,
		FromDate: fromDate,
		ToDate:   toDate,
	}

	posts, err := h.postService.Search(ctx, userID, filter)
	if err != nil {
		httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		return
	}

	out := make([]PostJSON, 0, len(posts))
	for _, p := range posts {
		out = append(out, mapPost(p))
	}

	resp := ListPostsResponse{
		Posts: out,
		Meta:  buildPageMeta(filter.Limit(), filter.Offset(), len(out)),
	}

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}
