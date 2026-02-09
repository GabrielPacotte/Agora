package userhttp

import (
	"net/http"
	"strconv"

	"github.com/GabrielPacotte/Agora/internal/filters"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func (h *Handler) ListSavedContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := httpcommon.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeAuthRequired, "unauthorized")
		return
	}

	limit, offset := parsePageQuery(r)

	page := filters.NewPage(limit, offset)
	limit = page.Limit()
	offset = page.Offset()

	list, err := h.userService.ListSavedContent(ctx, userID, page)
	if err != nil {
		httpcommon.WriteError(w, http.StatusInternalServerError, httpcommon.ErrCodeInternal, "internal server error")
		return
	}

	out := make([]SavedContentJSON, 0, len(list))
	for _, sc := range list {
		out = append(out, SavedContentJSON{
			UserID:      sc.UserID().String(),
			ContentID:   sc.SubjectID().String(),
			ContentType: string(sc.SubjectType()),
			SavedAt:     sc.SavedAt(),
		})
	}

	var nextOffset *int
	if len(list) == limit {
		n := offset + limit
		nextOffset = &n
	}

	resp := ListSavedContentResponse{
		Saved: out,
		Meta: PageMeta{
			NextOffset: nextOffset,
			Count:      len(out),
		},
	}

	httpcommon.WriteJSON(w, http.StatusOK, resp)
}

func parsePageQuery(r *http.Request) (limit, offset int) {
	limit = 0
	offset = 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	return limit, offset
}
